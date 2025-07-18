declare module 'vis-network-react' {
  import { Network, Data, Options } from 'vis-network'

  export interface GraphProps {
    graph: Data
    options?: Options
    getNetwork?: (network: Network) => void
    style?: React.CSSProperties
    className?: string
  }

  export const Graph: React.FC<GraphProps>
}
