// Libraries
import React, {PureComponent} from 'react'
import _ from 'lodash'

// Components
import {
  IndexList,
  IconFont,
  OverlayContainer,
  OverlayHeading,
  OverlayBody,
  Alignment,
  Greys,
  EmptyState,
  ComponentSize,
  Stack,
  ComponentSpacer,
} from 'src/clockface'

// Types
import {Authorization, Permission} from 'src/api'
const {ActionEnum, ResourceEnum} = Permission

interface Props {
  auth: Authorization
  onDismiss: () => void
}

export default class TokenRow extends PureComponent<Props> {
  public render() {
    const {auth, onDismiss} = this.props

    return (
      <OverlayContainer maxWidth={660}>
        <OverlayHeading title={auth.description} onDismiss={onDismiss} />
        <OverlayBody>
          <IndexList size={ComponentSize.Small}>
            <IndexList.Header>
              <IndexList.HeaderCell columnName="Resource" width="50%" />
              <IndexList.HeaderCell
                columnName="Read"
                width="25%"
                alignment={Alignment.Center}
              />
              <IndexList.HeaderCell
                columnName="Write"
                width="25%"
                alignment={Alignment.Center}
              />
            </IndexList.Header>
            <IndexList.Body emptyState={this.emptyList} columnCount={3}>
              {this.rows}
            </IndexList.Body>
          </IndexList>
        </OverlayBody>
      </OverlayContainer>
    )
  }

  private get rows(): JSX.Element[] {
    const {permissions} = this.props.auth

    return permissions.map(p => (
      <IndexList.Row>
        {this.resourceCell(p)}
        {this.readCell(p)}
        {this.writeCell(p)}
      </IndexList.Row>
    ))
  }

  private resourceCell = (permission: Permission): JSX.Element => {
    let resource = `All ${_.upperFirst(permission.resource)}`

    if (permission.id) {
      resource = permission.name
    }

    let resourceIcon: string

    switch (permission.resource) {
      case ResourceEnum.Buckets:
        resourceIcon = IconFont.Trash
        break
      case ResourceEnum.Users:
        resourceIcon = IconFont.User
        break
      case ResourceEnum.Orgs:
        resourceIcon = IconFont.UsersDuo
        break
      case ResourceEnum.Sources:
        resourceIcon = IconFont.Disks
        break
      case ResourceEnum.Dashboards:
        resourceIcon = IconFont.DashJ
        break
      default:
        resourceIcon = IconFont.Star
    }

    return (
      <IndexList.Cell>
        <ComponentSpacer stackChildren={Stack.Columns} align={Alignment.Left}>
          <span className={`icon ${resourceIcon}`} />
          <span>{resource}</span>
        </ComponentSpacer>
      </IndexList.Cell>
    )
  }

  private readCell = (permission: Permission): JSX.Element => {
    if (
      permission.action === ActionEnum.Read ||
      permission.action === ActionEnum.Write
    ) {
      return (
        <IndexList.Cell alignment={Alignment.Center}>
          <span
            className={`icon ${IconFont.Checkmark}`}
            style={{color: '#4ED8A0'}}
          />
        </IndexList.Cell>
      )
    }
  }

  private writeCell = (permission: Permission): JSX.Element => {
    if (permission.action === ActionEnum.Write) {
      return (
        <IndexList.Cell alignment={Alignment.Center}>
          <span
            className={`icon ${IconFont.Checkmark}`}
            style={{color: '#4ED8A0'}}
          />
        </IndexList.Cell>
      )
    }

    return (
      <IndexList.Cell alignment={Alignment.Center}>
        <span
          className={`icon ${IconFont.Remove}`}
          style={{color: `${Greys.Mountain}`}}
        />
      </IndexList.Cell>
    )
  }

  private get emptyList(): JSX.Element {
    return (
      <EmptyState>
        <EmptyState.Text text="This token has no associated permissions" />
      </EmptyState>
    )
  }
}
