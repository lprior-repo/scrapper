declare module 'vis-network-react' {
  import { Network, Data, Options } from 'vis-network'

  export interface GraphProps {
    readonly graph: Data
    readonly options?: Options
    readonly getNetwork?: (network: Network) => void
    readonly style?: React.CSSProperties
    readonly className?: string
  }

  export const Graph: React.FC<GraphProps>
}
