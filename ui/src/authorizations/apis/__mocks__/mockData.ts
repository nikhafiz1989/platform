import {Permission, Authorization} from 'src/api'

const {ActionEnum, ResourceEnum} = Permission
const {StatusEnum} = Authorization

export const authorization: Authorization = {
  links: {
    self: '/api/v2/authorizations/030444b11fb10000',
    user: '/api/v2/users/030444b10a710000',
  },
  org: 'Cooooool Orgggg bruhhh',
  description: 'I belive in free will',
  orgID: '0',
  id: '030444b11fb10000',
  token:
    'ohEmfY80A9UsW_cicNXgOMIPIsUvU6K9YcpTfCPQE3NV8Y6nTsCwVghczATBPyQh96CoZkOW5DIKldya6Y84KA==',
  status: StatusEnum.Active,
  user: 'watts',
  userID: '030444b10a710000',
  permissions: [
    {action: ActionEnum.Write, resource: ResourceEnum.Users},
    {action: ActionEnum.Write, resource: ResourceEnum.Orgs},
    {action: ActionEnum.Read, resource: ResourceEnum.Buckets},
    {
      action: ActionEnum.Write,
      resource: ResourceEnum.Buckets,
      id: '4',
      name: 'Bashful Bucket',
    },
    {
      action: ActionEnum.Write,
      resource: ResourceEnum.Buckets,
      id: '5',
      name: 'defbuck',
    },
  ],
}
