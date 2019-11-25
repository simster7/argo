import * as React from 'react';
import * as models from '../../../shared/models';

import * as classNames from 'classnames';
import {TagsInput} from '../../../shared/components/tags-input/tags-input';

require('./workflows-list.scss');

interface WorkflowsFilterProps {
    workflows: models.Workflow[];
}

interface WorkflowsFilterProps {}

export class WorkflowsFilter extends React.Component<WorkflowsFilterProps, {expanded: boolean}> {
    constructor(props: WorkflowsFilterProps) {
        super(props);
        this.state = {expanded: true};
    }

    public render() {
        return (
            <div className={classNames('workflows-list__filters-container', {'workflows-list__filters-container--expanded': this.state.expanded})}>
                <i
                    onClick={() => this.setState({expanded: !this.state.expanded})}
                    className={classNames('fa applications-list__filters-expander', {'fa-chevron-up': !this.state.expanded, 'fa-chevron-down': this.state.expanded})}
                />
                <p className='workflows-list__filters-container-title'>Filter By:</p>
                <div className='columns small-12 medium-3 xxlarge-12'>
                    <div className='workflows-list__filter'>
                        <p>Namespaces</p>
                        <ul>
                            <li>
                                <TagsInput
                                    // placeholder='*-us-west-*'
                                    // autocomplete={Array.from(new Set(applications.map(app => app.spec.destination.namespace).filter(item => !!item))).filter(
                                    //     ns => pref.namespacesFilter.indexOf(ns) === -1
                                    // )}
                                    tags={['simon', 'is', 'cool']}
                                    // onChange={selected => onChange({...pref, namespacesFilter: selected})}
                                />
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        );
    }
}
