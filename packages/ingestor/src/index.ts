import { Effect, Console, Layer, ConfigProvider } from 'effect'
import { NodeRuntime } from '@effect/platform-node'
import { GithubClient, GithubClientLive, Neo4jClientLive } from './services'
import {
  transformToGraphData,
  writeGraphData,
  parseCodeowners,
  type CodeownersEntry,
} from './logic'

// Main ingestion program
const program = Effect.gen(function* () {
  yield* Console.log('🚀 Starting GitHub organization ingestion...')

  // Get services
  const github = yield* GithubClient

  // Get organization name from command line or environment
  const orgName = process.argv[2] || process.env.GITHUB_ORG || 'golang'

  yield* Console.log(`📊 Fetching data for organization: ${orgName}`)

  // Fetch organization data
  const org = yield* github.fetchOrganization(orgName)
  yield* Console.log(`✅ Fetched organization: ${org.login}`)

  // Fetch repositories
  const repos = yield* github.fetchOrgRepos(orgName)
  yield* Console.log(`✅ Fetched ${repos.length} repositories`)

  // Fetch CODEOWNERS for each repository
  const codeownersMap = new Map<string, ReadonlyArray<CodeownersEntry>>()

  yield* Effect.forEach(
    repos,
    (repo) =>
      Effect.gen(function* () {
        const [owner, name] = repo.full_name.split('/')
        const codeownersContent = yield* github.fetchCodeowners(owner, name)

        if (codeownersContent) {
          const entries = parseCodeowners(codeownersContent)
          codeownersMap.set(repo.full_name, entries)
          yield* Console.log(
            `  ✅ Found CODEOWNERS for ${repo.full_name} (${entries.length} rules)`
          )
        }
      }).pipe(Effect.catchAll(() => Effect.succeed(undefined))),
    { concurrency: 5 }
  )

  yield* Console.log(
    `📋 Found CODEOWNERS files in ${codeownersMap.size}/${repos.length} repositories`
  )

  // Transform to graph data
  const graphData = yield* transformToGraphData(org, repos, codeownersMap)
  yield* Console.log(
    `🔄 Transformed data: ${graphData.nodes.length} nodes, ${graphData.edges.length} edges`
  )

  // Write to Neo4j
  yield* Console.log('💾 Writing to Neo4j database...')
  yield* writeGraphData(graphData)

  yield* Console.log('✨ Ingestion completed successfully!')
}).pipe(
  Effect.catchAll((error) =>
    Effect.gen(function* () {
      yield* Console.error('❌ Ingestion failed:')
      yield* Console.error(error)
      yield* Effect.fail(error)
    })
  )
)

// Load environment variables
import dotenv from 'dotenv'
dotenv.config({ path: '../../configs/.env' })

// Create config provider from environment
const configProvider = ConfigProvider.fromMap(
  new Map(Object.entries(process.env))
)

// Create the main layer with all dependencies
const MainLiveLayer = Layer.mergeAll(GithubClientLive, Neo4jClientLive).pipe(
  Layer.provide(
    Layer.mergeAll(NodeRuntime.layer, Layer.setConfigProvider(configProvider))
  )
)

// Run the program
program.pipe(Effect.provide(MainLiveLayer), NodeRuntime.runMain)
