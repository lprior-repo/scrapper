import React, { useState, useEffect } from 'react';
import { Effect } from 'effect';

import { CytoscapeGraphComponent } from './CytoscapeGraphComponent';
import { GraphErrorDisplay } from './GraphErrorDisplay';
import { GraphLoadingSpinner } from './GraphLoadingSpinner';
import type { GraphNode, GraphEdge } from '../services';

interface IGraphCanvasProps {
  readonly organization: string;
  readonly useTopics: boolean;
}

type GraphState =
  | { readonly type: 'loading' }
  | { readonly type: 'error'; readonly error: unknown }
  | { readonly type: 'success'; readonly data: { readonly nodes: readonly GraphNode[]; readonly edges: readonly GraphEdge[] } };

const createApiUrl = (organization: string, useTopics: boolean): string =>
  `http://localhost:8081/api/graph/${organization}${useTopics ? '?useTopics=true' : ''}`;

const fetchGraphData = (url: string) =>
  Effect.gen(function* () {
    const response = yield* Effect.tryPromise(() =>
      fetch(url).then((res) => {
        if (!res.ok) {
          throw new Error(`HTTP error! status: ${res.status}`);
        }
        return res.json();
      })
    );
    return response.data;
  });

const renderGraphState = (state: GraphState): React.ReactElement => {
  switch (state.type) {
    case 'loading':
      return <GraphLoadingSpinner />;
    case 'error':
      return <GraphErrorDisplay error={state.error} />;
    case 'success':
      return <CytoscapeGraphComponent nodes={state.data.nodes} edges={state.data.edges} />;
  }
};

export const GraphCanvas: React.FC<IGraphCanvasProps> = ({
  organization,
  useTopics,
}) => {
  const [state, setState] = useState<GraphState>({ type: 'loading' });

  useEffect(() => {
    const url = createApiUrl(organization, useTopics);
    const loadData = fetchGraphData(url);

    Effect.runPromise(loadData)
      .then((data) => {
        setState({ type: 'success', data });
      })
      .catch((error) => {
        setState({ type: 'error', error });
      });
  }, [organization, useTopics]);

  return renderGraphState(state);
};
