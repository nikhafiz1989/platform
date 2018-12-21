// Libraries
import React, {PureComponent, ChangeEvent} from 'react'

// Components
import {
  OverlayContainer,
  OverlayHeading,
  OverlayBody,
  Grid,
  Columns,
  Form,
  Button,
  ComponentColor,
  ComponentStatus,
} from 'src/clockface'
import TokenDescriptionInput from 'src/me/components/account/TokenDescriptionInput'
import TokenOrgDropdown from 'src/me/components/account/TokenOrgDropdown'
import TokenSelectedPermissions from 'src/me/components/account/TokenSelectedPermissions'

// Types
import {Authorization, Permission} from 'src/api'

interface Props {
  onGenerate: (authorization: Authorization) => void
  onDismiss: () => void
}

interface State {
  description: string
  orgID: string
  permissions: Permission[]
}

class GenerateTokenOverlay extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props)

    this.state = {
      description: '',
      orgID: '0',
      permissions: [],
    }
  }

  public render() {
    const {description, orgID, permissions} = this.state
    const {onDismiss} = this.props

    return (
      <OverlayContainer maxWidth={720}>
        <OverlayHeading title="Generate Token" onDismiss={onDismiss} />
        <OverlayBody>
          <Form onSubmit={this.handleFormSubmit}>
            <Grid>
              <Grid.Row>
                <Grid.Column widthSM={Columns.Six}>
                  <TokenDescriptionInput
                    description={description}
                    onChange={this.handleDescriptionChange}
                  />
                </Grid.Column>
                <Grid.Column widthSM={Columns.Six}>
                  <TokenOrgDropdown
                    selectedOrgID={orgID}
                    onChange={this.handleDropdownChange}
                  />
                </Grid.Column>
                <Grid.Column widthXS={Columns.Twelve}>
                  <TokenSelectedPermissions permissions={permissions} />
                </Grid.Column>
                <Grid.Column widthXS={Columns.Twelve}>
                  <Form.Footer>
                    <Button text="Cancel" onClick={onDismiss} />
                    <Button
                      text="Generate Token"
                      color={ComponentColor.Primary}
                      status={this.submitButtonStatus}
                    />
                  </Form.Footer>
                </Grid.Column>
              </Grid.Row>
            </Grid>
          </Form>
        </OverlayBody>
      </OverlayContainer>
    )
  }

  private handleDescriptionChange = (
    e: ChangeEvent<HTMLInputElement>
  ): void => {
    const description = e.target.value

    this.setState({description})
  }

  private handleDropdownChange = (orgID: string): void => {
    this.setState({orgID})
  }

  private handleFormSubmit = (): void => {
    console.log('generate token submission button clicked')
  }

  private get isFormValid(): boolean {
    const {description, orgID, permissions} = this.state

    if (!description || !orgID || !permissions.length) {
      return false
    }

    return true
  }

  private get submitButtonStatus(): ComponentStatus {
    if (this.isFormValid) {
      return ComponentStatus.Default
    }

    return ComponentStatus.Disabled
  }
}

export default GenerateTokenOverlay
