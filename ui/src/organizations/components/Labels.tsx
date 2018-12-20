// Libraries
import React, {PureComponent, ChangeEvent} from 'react'

// Components
import CreateLabelOverlay from 'src/organizations/components/CreateLabelOverlay'
import ProfilePageHeader from 'src/shared/components/profile_page/ProfilePageHeader'
import {
  ComponentSize,
  EmptyState,
  IconFont,
  Input,
  Button,
  ComponentColor,
} from 'src/clockface'
import LabelList from 'src/organizations/components/LabelList'
import FilterList from 'src/shared/components/Filter'

// API
import {createLabel, deleteLabel, updateLabel} from 'src/organizations/apis'

// Types
import {LabelType} from 'src/clockface'
import {Label, Organization} from 'src/api'

// Decorators
import {ErrorHandling} from 'src/shared/decorators/errors'

interface Props {
  labels: Label[]
  org: Organization
}

interface State {
  searchTerm: string
  isOverlayVisible: boolean
  labelTypes: LabelType[]
}

@ErrorHandling
export default class Labels extends PureComponent<Props, State> {
  constructor(props) {
    super(props)

    this.state = {
      searchTerm: '',
      isOverlayVisible: false,
      labelTypes: this.labelTypes(this.props.labels),
    }
  }

  public render() {
    const {org} = this.props
    const {searchTerm, isOverlayVisible, labelTypes} = this.state

    return (
      <>
        <ProfilePageHeader>
          <Input
            icon={IconFont.Search}
            widthPixels={290}
            value={searchTerm}
            onBlur={this.handleFilterBlur}
            onChange={this.handleFilterChange}
            placeholder="Filter Labels..."
          />
          <Button
            text="Create Labels"
            color={ComponentColor.Primary}
            icon={IconFont.Plus}
            onClick={this.handleShowOverlay}
          />
        </ProfilePageHeader>
        <FilterList<LabelType>
          list={labelTypes}
          searchKeys={['name', 'description']}
          searchTerm={searchTerm}
        >
          {ls => (
            <LabelList
              labels={ls}
              emptyState={this.emptyState}
              onUpdateLabel={this.handleUpdateLabel}
            />
          )}
        </FilterList>
        <CreateLabelOverlay
          org={org}
          isVisible={isOverlayVisible}
          onDismiss={this.handleDismissOverlay}
          onCreateLabel={this.handleCreateLabel}
        />
      </>
    )
  }

  private handleShowOverlay = (): void => {
    this.setState({isOverlayVisible: true})
  }

  private handleDismissOverlay = (): void => {
    this.setState({isOverlayVisible: false})
  }

  private handleFilterChange = (e: ChangeEvent<HTMLInputElement>): void => {
    this.setState({searchTerm: e.target.value})
  }

  private handleFilterBlur = (e: ChangeEvent<HTMLInputElement>): void => {
    this.setState({searchTerm: e.target.value})
  }

  private handleCreateLabel = async (
    org: Organization,
    labelType: LabelType
  ) => {
    const newLabel = await createLabel(org, {
      name: labelType.name,
      color: labelType.colorHex,
      description: labelType.description,
    })
    const labelTypes = this.labelTypes([...this.props.labels, newLabel])
    this.setState({labelTypes})
  }

  private handleUpdateLabel = async (labelType: LabelType) => {
    await updateLabel(this.props.org, {
      name: labelType.name,
      color: labelType.colorHex,
      description: labelType.description,
    })

    const labelTypes = this.state.labelTypes.map(l => {
      if (l.id === labelType.id) {
        return labelType
      }

      return l
    })

    this.setState({labelTypes})
  }

  private labelTypes(labels: Label[]): LabelType[] {
    return labels.map(label => ({
      id: label.name,
      name: label.name,
      description: label.description,
      colorHex: label.color,
      onDelete: this.handleDelete,
    }))
  }

  private handleDelete = async (name: string) => {
    const {org, labels} = this.props
    const label = labels.find(label => label.name === name)

    await deleteLabel(org, label)
    const labelTypes = this.state.labelTypes.filter(l => l.id !== name)

    this.setState({labelTypes})
  }

  private get emptyState(): JSX.Element {
    const {searchTerm} = this.state

    let emptyText = "This organization doesn't have any Labels yet"

    if (searchTerm) {
      emptyText = 'No Labels match your search term'
    }

    return (
      <EmptyState size={ComponentSize.Medium}>
        <EmptyState.Text text={emptyText} />
      </EmptyState>
    )
  }
}
