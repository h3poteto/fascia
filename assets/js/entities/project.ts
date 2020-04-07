export type ServerProject = {
  ID: number
  UserID: number
  Title: string
  Description: string
  RepositoryID: number | null
  ShowIssues: boolean
  ShowPullRequests: boolean
}

export type Project = {
  id: number
  userID: number
  title: string
  description: string
  repositoryID: number | null
  showIssues: boolean
  showPullRequests: boolean
}

export const converter = (p: ServerProject): Project => ({
  id: p.ID,
  userID: p.UserID,
  title: p.Title,
  description: p.Description,
  repositoryID: p.RepositoryID,
  showIssues: p.ShowIssues,
  showPullRequests: p.ShowPullRequests
})
