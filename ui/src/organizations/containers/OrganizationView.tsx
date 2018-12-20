// Libraries
import React, {PureComponent} from 'react'
import {WithRouterProps} from 'react-router'
import {connect} from 'react-redux'
import _ from 'lodash'

// APIs
import {
  getBuckets,
  getDashboards,
  getMembers,
  getTasks,
  getLabels,
} from 'src/organizations/apis'

// Actions
import {updateOrg} from 'src/organizations/actions'

// Components
import {Page} from 'src/pageLayout'
import {Spinner} from 'src/clockface'
import ProfilePage from 'src/shared/components/profile_page/ProfilePage'
import ProfilePageSection from 'src/shared/components/profile_page/ProfilePageSection'
import Members from 'src/organizations/components/Members'
import Buckets from 'src/organizations/components/Buckets'
import Dashboards from 'src/organizations/components/Dashboards'
import Tasks from 'src/organizations/components/Tasks'
import Labels from 'src/organizations/components/Labels'
import OrgOptions from 'src/organizations/components/OrgOptions'
import GetOrgResources from 'src/organizations/components/GetOrgResources'

// Types
import {AppState, Dashboard} from 'src/types/v2'
import {LabelType} from 'src/clockface'
import {User, Bucket, Organization, Task} from 'src/api'

// Decorators
import {ErrorHandling} from 'src/shared/decorators/errors'

interface StateProps {
  org: Organization
}

interface DispatchProps {
  onUpdateOrg: typeof updateOrg
}

type Props = StateProps & WithRouterProps & DispatchProps

@ErrorHandling
class OrganizationView extends PureComponent<Props> {
  public render() {
    const {org, params, onUpdateOrg} = this.props

    return (
      <Page>
        <Page.Header fullWidth={false}>
          <Page.Header.Left>
            <Page.Title title={org.name ? org.name : 'Organization'} />
          </Page.Header.Left>
          <Page.Header.Right />
        </Page.Header>
        <Page.Contents fullWidth={false} scrollable={true}>
          <div className="col-xs-12">
            <ProfilePage
              name={org.name}
              parentUrl={`/organizations/${org.id}`}
              activeTabUrl={params.tab}
            >
              <ProfilePageSection
                id="org-view-tab--members"
                url="members_tab"
                title="Members"
              >
                <GetOrgResources<User[]>
                  organization={org}
                  fetcher={getMembers}
                >
                  {(members, loading) => (
                    <Spinner loading={loading}>
                      <Members members={members} />
                    </Spinner>
                  )}
                </GetOrgResources>
              </ProfilePageSection>
              <ProfilePageSection
                id="org-view-tab--buckets"
                url="buckets_tab"
                title="Buckets"
              >
                <GetOrgResources<Bucket[]>
                  organization={org}
                  fetcher={getBuckets}
                >
                  {(buckets, loading) => (
                    <Spinner loading={loading}>
                      <Buckets buckets={buckets} org={org} />
                    </Spinner>
                  )}
                </GetOrgResources>
              </ProfilePageSection>
              <ProfilePageSection
                id="org-view-tab--dashboards"
                url="dashboards_tab"
                title="Dashboards"
              >
                <GetOrgResources<Dashboard[]>
                  organization={org}
                  fetcher={getDashboards}
                >
                  {(dashboards, loading) => (
                    <Spinner loading={loading}>
                      <Dashboards dashboards={dashboards} />
                    </Spinner>
                  )}
                </GetOrgResources>
              </ProfilePageSection>
              <ProfilePageSection
                id="org-view-tab--tasks"
                url="tasks_tab"
                title="Tasks"
              >
                <GetOrgResources<Task[]> organization={org} fetcher={getTasks}>
                  {(tasks, loading) => (
                    <Spinner loading={loading}>
                      <Tasks tasks={tasks} />
                    </Spinner>
                  )}
                </GetOrgResources>
              </ProfilePageSection>
              <ProfilePageSection
                id="org-view-tab--labels"
                url="labels_tab"
                title="Labels"
              >
                <GetOrgResources<LabelType[]>
                  organization={org}
                  fetcher={getLabels}
                >
                  {(labels, loading) => (
                    <Spinner loading={loading}>
                      <Labels labels={labels} org={org} />
                    </Spinner>
                  )}
                </GetOrgResources>
              </ProfilePageSection>
              <ProfilePageSection
                id="org-view-tab--options"
                url="options_tab"
                title="Options"
              >
                <OrgOptions org={org} onUpdateOrg={onUpdateOrg} />
              </ProfilePageSection>
            </ProfilePage>
          </div>
        </Page.Contents>
      </Page>
    )
  }
}

const mstp = (state: AppState, props: WithRouterProps) => {
  const {orgs} = state
  const org = orgs.find(o => o.id === props.params.orgID)
  return {
    org,
  }
}

const mdtp: DispatchProps = {
  onUpdateOrg: updateOrg,
}

export default connect<StateProps, DispatchProps, {}>(
  mstp,
  mdtp
)(OrganizationView)
