import { Config } from 'effect'

export const Neo4jConfig = Config.all({
  uri: Config.string('NEO4J_URI').pipe(
    Config.withDefault('bolt://localhost:7687')
  ),
  username: Config.string('NEO4J_USERNAME').pipe(Config.withDefault('neo4j')),
  password: Config.string('NEO4J_PASSWORD').pipe(
    Config.withDefault('password')
  ),
  database: Config.string('NEO4J_DATABASE').pipe(Config.withDefault('neo4j')),
})

export const GithubConfig = Config.all({
  token: Config.string('GITHUB_TOKEN'),
  baseUrl: Config.string('GITHUB_BASE_URL').pipe(
    Config.withDefault('https://api.github.com')
  ),
  userAgent: Config.string('GITHUB_USER_AGENT').pipe(
    Config.withDefault('overseer-codeowners-scanner/1.0')
  ),
  maxRepos: Config.integer('GITHUB_MAX_REPOS').pipe(Config.withDefault(100)),
  maxTeams: Config.integer('GITHUB_MAX_TEAMS').pipe(Config.withDefault(50)),
})

export const AppConfig = Config.all({
  neo4j: Neo4jConfig,
  github: GithubConfig,
})
