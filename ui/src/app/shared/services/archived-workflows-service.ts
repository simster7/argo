import * as models from '../../../models';
import requests from './requests';

export class ArchivedWorkflowsService {
    public list() {
        return requests
            .get(`api/v1/archived-workflows/`)
            .then(res => res.body as models.WorkflowList)
            .then(list => list.items || []);
    }

    public get(namespace: string, uid: string): Promise<models.Workflow> {
        return requests.get(`api/v1/archived-workflows/${namespace}/${uid}`).then(res => res.body as models.Workflow);
    }

    public resubmit(namespace: string, uid: string) {
        return requests.put(`api/v1/archived-workflows/${namespace}/${uid}/resubmit`).then(res => res.body as models.Workflow);
    }

    public delete(namespace: string, uid: string) {
        return requests.delete(`api/v1/archived-workflows/${namespace}/${uid}`);
    }
}
