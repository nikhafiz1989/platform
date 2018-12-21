// Libraries
import React, {PureComponent} from 'react'

// Components
import {
  IndexList,
  ComponentSpacer,
  Alignment,
  Button,
  ComponentSize,
  ComponentColor,
  OverlayTechnology,
} from 'src/clockface'
import ViewTokenOverlay from 'src/me/components/account/ViewTokenOverlay'

// Types
import {Authorization} from 'src/api'
import {OverlayState} from 'src/types/v2'

// Actions
import {notify} from 'src/shared/actions/notifications'

interface Props {
  auth: Authorization
  onNotify: typeof notify
  onDelete: (authID: string) => void
}

interface State {
  overlayState: OverlayState
}

export default class TokenRow extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props)

    this.state = {
      overlayState: OverlayState.Closed,
    }
  }

  public render() {
    const {auth} = this.props
    const {description, status, org, orgID} = auth

    return (
      <>
        <IndexList.Row>
          <IndexList.Cell>
            <a href="#" onClick={this.handleShowOverlay}>
              {description}
            </a>
          </IndexList.Cell>
          <IndexList.Cell>{status}</IndexList.Cell>
          <IndexList.Cell>{org}</IndexList.Cell>
          <IndexList.Cell alignment={Alignment.Right} revealOnHover={true}>
            <ComponentSpacer align={Alignment.Right}>
              <Button
                size={ComponentSize.ExtraSmall}
                color={ComponentColor.Danger}
                text="Delete"
                onClick={this.handleDelete}
              />
            </ComponentSpacer>
          </IndexList.Cell>
        </IndexList.Row>
        <OverlayTechnology visible={this.overlayVisibility}>
          <ViewTokenOverlay auth={auth} onDismiss={this.handleDismissOverlay} />
        </OverlayTechnology>
      </>
    )
  }

  private get overlayVisibility(): boolean {
    const {overlayState} = this.state

    return overlayState === OverlayState.Open
  }

  private handleShowOverlay = (): void => {
    this.setState({overlayState: OverlayState.Open})
  }

  private handleDismissOverlay = (): void => {
    this.setState({overlayState: OverlayState.Closed})
  }

  private handleDelete = () => {
    const {
      auth: {id},
      onDelete,
    } = this.props
    onDelete(id)
  }
}
