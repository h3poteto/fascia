import expect from 'expect'
import ProjectReducer from '../../frontend/javascripts/reducers/ProjectReducer'
import * as projectActions from '../../frontend/javascripts/actions/ProjectAction'
import 'babel-polyfill'

// shared examples
function sharedExampleInitState(action) {
  expect(
    ProjectReducer(undefined, action)
  ).toEqual({
    isModalOpen: false,
    newProject: {title: "", description: ""},
    projects: [],
    repositories: [],
    selectedRepository: null,
    isLoading: false,
    error: null
  })
}

describe('ProjectReducer', () => {
  describe('initState', () => {
    it('should return the initial state', () => {
      sharedExampleInitState({})
    })
  })
  describe('SERVER_ERROR', () => {
    it('should return server error', () => {
      expect(
        ProjectReducer(null, {
          type: projectActions.SERVER_ERROR
        })
      ).toEqual({
        error: "Internal Server Error",
        isLoading: false
      })
    })
  })
  describe('CLOSE_FLASH', () => {
    it('should close flash', () => {
      expect(
        ProjectReducer({
          error: "Internal Server Error"
        }, {
          type: projectActions.CLOSE_FLASH
        })
      ).toEqual({
        error: null
      })
    })
  })
  describe('REQUEST_POSTS', () => {
    it('should return previous state', () => {
      sharedExampleInitState({type: projectActions.REQUEST_POSTS})
    })
  })
  describe('RECEIVE_POSTS', () => {
    context('when projects are empty', ()=> {
      it('should return empty projects', () => {
        expect(
          ProjectReducer({
            projects: []
          }, {
            type: projectActions.RECEIVE_POSTS,
            projects: null
          })
        ).toEqual({
          projects: []
        })
      })
    })
    context('when projects are not empty', ()=> {
      it('should return projects', () => {
        expect(
          ProjectReducer({
            projects: []
          }, {
            type: projectActions.RECEIVE_POSTS,
            projects: ["project1", "project2"]
          })
        ).toEqual({
          projects: ["project1", "project2"]
        })
      })
    })
  })
  describe('OPEN_NEW_PROJECT', () => {
    it('should open project modal', () => {
      expect(
        ProjectReducer({
          isModalOpen: false
        }, {
          type: projectActions.OPEN_NEW_PROJECT,
          isModalOpen: true
        })
      ).toEqual({
        isModalOpen: true
      })
    })
  })
  describe('CLOSE_NEW_PROJECT', () => {
    it('should close project modal', () => {
      expect(
        ProjectReducer({
          isModalOpen: true
        }, {
          type: projectActions.CLOSE_NEW_PROJECT,
          isModalOpen: false
        })
      ).toEqual({
        isModalOpen: false
      })
    })
  })

  describe('REQUEST_REPOSITORIES', () => {
    it('should return previous state', () => {
      sharedExampleInitState({ type: projectActions.REQUEST_REPOSITORIES })
    })
  })

  describe('RECEIVE_REPOSITORIES', () => {
    it('should return repositories', () => {
      expect(
        ProjectReducer({
          repositories: []
        }, {
          type: projectActions.RECEIVE_REPOSITORIES,
          repositories: ["repo1", "repo2"]
        })
      ).toEqual({
        repositories: ["repo1", "repo2"]
      })
    })
  })

  describe('CHANGE_SELECT_REPOSITORY', () => {
    it('should return repository object and set new project title', () => {
      const stateRepositories = [
        {
          id: 2,
          title: "repo1"
        },{
          id: 3,
          title: "repo2"
        }
      ]
      expect(
        ProjectReducer({
          repositories: stateRepositories,
          newProject: {title: "", description: ""}
        }, {
          type: projectActions.CHANGE_SELECT_REPOSITORY,
          selectEvent: {
            selectedIndex: 0,
            options: [
              { text: "repo1" },
              { text: "repo2" }
            ],
            value: 2
          }
        })
      ).toEqual({
        repositories: stateRepositories,
        selectedRepository: {
          id: 2,
          title: "repo1"
        },
        newProject: {title: "repo1", description: ""}
      })
    })
  })

  describe('REQUEST_CREATE_PROJECT', () => {
    it('should close modal and open whole loading', () => {
      expect(
        ProjectReducer({
          projects: ["project1", "project2"],
          newProject: { title: "project3", description: "" },
          isModalOpen: true,
          isLoading: false
        }, {
          type:projectActions.REQUEST_CREATE_PROJECT,
          project: "project3"
        })
      ).toEqual({
        projects: ["project1", "project2"],
        newProject: { title: "project3", description: "" },
        isModalOpen: false,
        isLoading: true
      })
    })
  })

  describe('RECEIVE_CREATE_PROJECT', () => {
    it('should return projects', () => {
      expect(
        ProjectReducer({
          projects: ["project1", "project2"],
          newProject: { title: "project3", description: "" },
          isLoading: true
        },{
          type: projectActions.RECEIVE_CREATE_PROJECT,
          project: "project3"
        })
      ).toEqual({
        projects: ["project1", "project2", "project3"],
        newProject: {title: "", description: ""},
        isLoading: false
      })
    })
  })

  describe('UPDATE_NEW_PROJECT_TITLE', () => {
    it('should return new project title', () => {
      expect(
        ProjectReducer({
          newProject: { title: "pro", description: "" }
        },{
          type: projectActions.UPDATE_NEW_PROJECT_TITLE,
          title: "proj"
        })
      ).toEqual({
        newProject: { title: "proj", description: "" }
      })
    })
  })

  describe('UPDATE_NEW_PROJECT_DESCRIPTION', () => {
    it('should return new project description', () => {
      expect(
        ProjectReducer({
          newProject: { title: "project1", description: "project des" }
        },{
          type: projectActions.UPDATE_NEW_PROJECT_DESCRIPTION,
          description: "project description"
        })
      ).toEqual({
        newProject: { title: "project1", description: "project description" }
      })
    })
  })
})
