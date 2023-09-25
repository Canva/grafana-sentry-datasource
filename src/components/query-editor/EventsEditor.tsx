import React from 'react';
import { InlineFormLabel, Input, Select } from '@grafana/ui';
import { SentryDataSource } from '../../datasource';
import { selectors } from '../../selectors';
import { SentryEventSortOptions } from '../../constants';
import type { QueryEditorProps } from '@grafana/data';
import type { SentryConfig, SentryQuery, SentryEventSort } from '../../types';

type EventsEditorProps = Pick<QueryEditorProps<SentryDataSource, SentryQuery, SentryConfig>, 'query' | 'onChange' | 'onRunQuery'>;

export const EventsEditor = ({ query, onChange, onRunQuery }: EventsEditorProps) => {
  const onEventsQueryChange = (eventsQuery: string) => {
    onChange({ ...query, eventsQuery } as SentryQuery);
  };
  const onEventsSortChange = (eventsSort: SentryEventSort) => {
    onChange({ ...query, eventsSort } as SentryQuery);
    onRunQuery();
  };
  const onEventsLimitChange = (eventsLimit?: number) => {
    onChange({ ...query, eventsLimit } as SentryQuery);
  };
  return query.queryType === 'events' ? (
    <>
      <div className="gf-form">
        <InlineFormLabel width={10} className="query-keyword" tooltip={selectors.components.QueryEditor.Events.Query.tooltip}>
          {selectors.components.QueryEditor.Events.Query.label}
        </InlineFormLabel>
        <Input
          value={query.eventsQuery}
          onChange={(e) => onEventsQueryChange(e.currentTarget.value)}
          onBlur={onRunQuery}
          placeholder={selectors.components.QueryEditor.Events.Query.placeholder}
        />
      </div>
      <div className="gf-form">
        <InlineFormLabel width={10} className="query-keyword" tooltip={selectors.components.QueryEditor.Events.Sort.tooltip}>
          {selectors.components.QueryEditor.Events.Sort.label}
        </InlineFormLabel>
        <Select
          options={SentryEventSortOptions}
          value={query.eventsSort}
          width={28}
          onChange={(e) => onEventsSortChange(e?.value!)}
          className="inline-element"
          placeholder={selectors.components.QueryEditor.Events.Sort.placeholder}
          isClearable={true}
        />
        <InlineFormLabel width={8} className="query-keyword" tooltip={selectors.components.QueryEditor.Events.Limit.tooltip}>
          {selectors.components.QueryEditor.Events.Limit.label}
        </InlineFormLabel>
        <Input
          value={query.eventsLimit}
          type="number"
          onChange={(e) => onEventsLimitChange(e.currentTarget.valueAsNumber)}
          onBlur={onRunQuery}
          width={32}
          className="inline-element"
          placeholder={selectors.components.QueryEditor.Events.Limit.placeholder}
        />
      </div>
    </>
  ) : (
    <></>
  );
};
