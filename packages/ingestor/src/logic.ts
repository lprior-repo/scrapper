import { Effect } from 'effect'
import { Schema } from '@effect/schema'
import type {
  GitHubRepository,
  GitHubOrganization,
  Neo4jClient,
} from './services'

// Graph Data Types
export const GraphNode = Schema.Struct({
  id: Schema.String,
  type: Schema.Literal('organization', 'repository', 'user', 'team'),
  properties: Schema.Record(Schema.String, Schema.Unknown),
})

export type GraphNode = Schema.Schema.Type<typeof GraphNode>

export const GraphEdge = Schema.Struct({
  from: Schema.String,
  to: Schema.String,
  type: Schema.Literal('OWNS', 'HAS_CODEOWNER', 'MEMBER_OF'),
  properties: Schema.Record(Schema.String, Schema.Unknown),
})

export type GraphEdge = Schema.Schema.Type<typeof GraphEdge>

export const GraphData = Schema.Struct({
  nodes: Schema.Array(GraphNode),
  edges: Schema.Array(GraphEdge),
})

export type GraphData = Schema.Schema.Type<typeof GraphData>

// CODEOWNERS Parsing
export interface CodeownersEntry {
  pattern: string
  owners: ReadonlyArray<string>
  line: number
}

export const parseCodeowners = (
  content: string
): ReadonlyArray<CodeownersEntry> => {
  const lines = content.split('\n')
  const entries: CodeownersEntry[] = []

  lines.forEach((line, index) => {
    const trimmed = line.trim()

    // Skip comments and empty lines
    if (!trimmed || trimmed.startsWith('#')) return

    const parts = trimmed.split(/\s+/)
    if (parts.length < 2) return

    const pattern = parts[0]
    const owners = parts
      .slice(1)
      .filter((owner) => owner.startsWith('@') || owner.includes('/'))

    if (owners.length > 0) {
      entries.push({
        pattern,
        owners,
        line: index + 1,
      })
    }
  })

  return entries
}

// Transform Functions
export const transformToGraphData = (
  org: GitHubOrganization,
  repos: ReadonlyArray<GitHubRepository>,
  codeownersMap: Map<string, ReadonlyArray<CodeownersEntry>>
): Effect.Effect<GraphData, never> =>
  Effect.succeed({
    nodes: [
      // Organization node
      {
        id: `org:${org.login}`,
        type: 'organization' as const,
        properties: {
          login: org.login,
          name: org.name || org.login,
          description: org.description || '',
          created_at: org.created_at,
          updated_at: org.updated_at,
        },
      },
      // Repository nodes
      ...repos.map((repo) => ({
        id: `repo:${repo.full_name}`,
        type: 'repository' as const,
        properties: {
          name: repo.name,
          full_name: repo.full_name,
          description: repo.description || '',
          url: repo.url,
          private: repo.private,
          default_branch: repo.default_branch,
          created_at: repo.created_at,
          updated_at: repo.updated_at,
          has_codeowners: codeownersMap.has(repo.full_name),
        },
      })),
      // Extract unique owners as nodes
      ...Array.from(
        new Set(
          Array.from(codeownersMap.values())
            .flatMap((entries) => entries)
            .flatMap((entry) => entry.owners)
        )
      ).map((owner) => ({
        id:
          owner.startsWith('@') && !owner.includes('/')
            ? `user:${owner.substring(1)}`
            : `team:${owner}`,
        type: (owner.startsWith('@') && !owner.includes('/')
          ? 'user'
          : 'team') as const,
        properties: {
          name: owner,
          login: owner.startsWith('@') ? owner.substring(1) : owner,
        },
      })),
    ],
    edges: [
      // Organization owns repositories
      ...repos.map((repo) => ({
        from: `org:${org.login}`,
        to: `repo:${repo.full_name}`,
        type: 'OWNS' as const,
        properties: {},
      })),
      // Repositories have codeowners
      ...Array.from(codeownersMap.entries()).flatMap(
        ([repoFullName, entries]) =>
          entries.flatMap((entry) =>
            entry.owners.map((owner) => ({
              from: `repo:${repoFullName}`,
              to:
                owner.startsWith('@') && !owner.includes('/')
                  ? `user:${owner.substring(1)}`
                  : `team:${owner}`,
              type: 'HAS_CODEOWNER' as const,
              properties: {
                pattern: entry.pattern,
                line: entry.line,
              },
            }))
          )
      ),
    ],
  })

// Database Writing
export const writeGraphData = (
  data: GraphData
): Effect.Effect<void, Error, Neo4jClient> =>
  Effect.gen(function* () {
    const neo4j = yield* Neo4jClient

    // Clear existing data
    yield* neo4j.run('MATCH (n) DETACH DELETE n')

    // Create constraints
    const constraints = [
      'CREATE CONSTRAINT IF NOT EXISTS FOR (o:Organization) REQUIRE o.login IS UNIQUE',
      'CREATE CONSTRAINT IF NOT EXISTS FOR (r:Repository) REQUIRE r.full_name IS UNIQUE',
      'CREATE CONSTRAINT IF NOT EXISTS FOR (u:User) REQUIRE u.login IS UNIQUE',
      'CREATE CONSTRAINT IF NOT EXISTS FOR (t:Team) REQUIRE t.name IS UNIQUE',
    ]

    yield* Effect.forEach(constraints, (query) => neo4j.run(query), {
      concurrency: 1,
    })

    // Create nodes
    yield* Effect.forEach(
      data.nodes,
      (node) => {
        const label = node.type.charAt(0).toUpperCase() + node.type.slice(1)
        const query = `
          CREATE (n:${label} $props)
          SET n.id = $id
          RETURN n
        `
        return neo4j.run(query, { id: node.id, props: node.properties })
      },
      { concurrency: 'unbounded' }
    )

    // Create edges
    yield* Effect.forEach(
      data.edges,
      (edge) => {
        const query = `
          MATCH (from {id: $fromId})
          MATCH (to {id: $toId})
          CREATE (from)-[r:${edge.type} $props]->(to)
          RETURN r
        `
        return neo4j.run(query, {
          fromId: edge.from,
          toId: edge.to,
          props: edge.properties,
        })
      },
      { concurrency: 'unbounded' }
    )

    console.log(
      `✅ Created ${data.nodes.length} nodes and ${data.edges.length} edges`
    )
  })
