import { Context, Effect, Layer, Scope, Either } from 'effect'
import {
  HttpClient,
  HttpClientError,
  HttpClientRequest,
  HttpClientResponse,
} from '@effect/platform'
import neo4j, { Session, QueryResult, Record } from 'neo4j-driver'
import { AppConfig } from './config'
import { Schema, ParseResult } from '@effect/schema'

// GitHub Types
export const GitHubRepository = Schema.Struct({
  id: Schema.Number,
  name: Schema.String,
  full_name: Schema.String,
  description: Schema.NullOr(Schema.String),
  url: Schema.String,
  private: Schema.Boolean,
  created_at: Schema.String,
  updated_at: Schema.String,
  default_branch: Schema.String,
})

export type GitHubRepository = Schema.Schema.Type<typeof GitHubRepository>

export const GitHubOrganization = Schema.Struct({
  id: Schema.Number,
  login: Schema.String,
  name: Schema.NullOr(Schema.String),
  description: Schema.NullOr(Schema.String),
  url: Schema.String,
  created_at: Schema.String,
  updated_at: Schema.String,
})

export type GitHubOrganization = Schema.Schema.Type<typeof GitHubOrganization>

// Service Interfaces
export interface GithubClient {
  readonly fetchOrgRepos: (
    org: string
  ) => Effect.Effect<
    ReadonlyArray<GitHubRepository>,
    HttpClientError | ParseResult.ParseError
  >
  readonly fetchOrganization: (
    org: string
  ) => Effect.Effect<
    GitHubOrganization,
    HttpClientError | ParseResult.ParseError
  >
  readonly fetchCodeowners: (
    owner: string,
    repo: string
  ) => Effect.Effect<string | null, HttpClientError>
}

export const GithubClient = Context.GenericTag<GithubClient>('GithubClient')

export interface Neo4jClient {
  readonly run: (
    query: string,
    params?: Record<string, unknown>
  ) => Effect.Effect<QueryResult, Error>
  readonly runInTransaction: <T>(
    fn: (tx: Session) => Promise<T>
  ) => Effect.Effect<T, Error>
}

export const Neo4jClient = Context.GenericTag<Neo4jClient>('Neo4jClient')

// Live Service Implementations
export const GithubClientLive = Layer.effect(
  GithubClient,
  Effect.gen(function* () {
    const config = yield* AppConfig
    const httpClient = yield* HttpClient.HttpClient

    const defaultHeaders = {
      Accept: 'application/vnd.github.v3+json',
      'User-Agent': config.github.userAgent,
      Authorization: `token ${config.github.token}`,
    }

    const fetchOrgRepos = (org: string) =>
      Effect.gen(function* () {
        const url = `${config.github.baseUrl}/orgs/${org}/repos`
        const request = HttpClientRequest.get(url).pipe(
          HttpClientRequest.setHeaders(defaultHeaders)
        )

        const response = yield* httpClient(request)
        const json = yield* HttpClientResponse.json(response)
        const repos = yield* Schema.decodeUnknown(
          Schema.Array(GitHubRepository)
        )(json)

        return repos.slice(0, config.github.maxRepos)
      })

    const fetchOrganization = (org: string) =>
      Effect.gen(function* () {
        const url = `${config.github.baseUrl}/orgs/${org}`
        const request = HttpClientRequest.get(url).pipe(
          HttpClientRequest.setHeaders(defaultHeaders)
        )

        const response = yield* httpClient(request)
        const json = yield* HttpClientResponse.json(response)
        return yield* Schema.decodeUnknown(GitHubOrganization)(json)
      })

    const fetchCodeowners = (owner: string, repo: string) =>
      Effect.gen(function* () {
        const locations = [
          '.github/CODEOWNERS',
          'CODEOWNERS',
          'docs/CODEOWNERS',
        ]

        for (const location of locations) {
          const url = `${config.github.baseUrl}/repos/${owner}/${repo}/contents/${location}`
          const request = HttpClientRequest.get(url).pipe(
            HttpClientRequest.setHeaders(defaultHeaders)
          )

          const result = yield* Effect.either(
            httpClient(request).pipe(
              Effect.flatMap(HttpClientResponse.json),
              Effect.flatMap((json: { content?: string; type?: string }) =>
                Effect.try(() => {
                  if (json.content) {
                    return Buffer.from(json.content, 'base64').toString('utf-8')
                  }
                  return null
                })
              )
            )
          )

          if (Either.isRight(result) && result.right) {
            return result.right
          }
        }

        return null
      })

    return GithubClient.of({
      fetchOrgRepos,
      fetchOrganization,
      fetchCodeowners,
    })
  }).pipe(Layer.provide(HttpClient.layer))
)

export const Neo4jClientLive = Layer.scoped(
  Neo4jClient,
  Effect.gen(function* () {
    const config = yield* AppConfig
    yield* Scope.Scope

    const driver = yield* Effect.acquireRelease(
      Effect.try(() =>
        neo4j.driver(
          config.neo4j.uri,
          neo4j.auth.basic(config.neo4j.username, config.neo4j.password)
        )
      ),
      (driver) => Effect.promise(() => driver.close())
    )

    const run = (query: string, params?: Record<string, unknown>) =>
      Effect.gen(function* () {
        const session = driver.session({
          database: config.neo4j.database,
        })

        try {
          const result = yield* Effect.promise(() => session.run(query, params))
          return result
        } finally {
          yield* Effect.promise(() => session.close())
        }
      })

    const runInTransaction = <T>(fn: (tx: Session) => Promise<T>) =>
      Effect.gen(function* () {
        const session = driver.session({
          database: config.neo4j.database,
        })

        try {
          return yield* Effect.promise(() =>
            session.executeWrite((tx) => fn(tx as Session))
          )
        } finally {
          yield* Effect.promise(() => session.close())
        }
      })

    return Neo4jClient.of({
      run,
      runInTransaction,
    })
  })
)
