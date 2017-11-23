import expect from 'expect'
import ProjectReducer from '../../frontend/javascripts/reducers/ProjectReducer'
import * as projectActions from '../../frontend/javascripts/actions/ProjectAction'
import * as newProjectModalActions from '../../frontend/javascripts/actions/ProjectAction/NewProjectModalAction'
import 'babel-polyfill'

// shared examples
function sharedExampleInitState(action) {
  expect(
    ProjectReducer(undefined, action)
  ).toEqual({
    isModalOpen: false,
    projects: [],
    repositories: [],
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

  context('projectActions', () => {
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
  })

  context('newProjectModalActions', () => {
    describe('CLOSE_NEW_PROJECT', () => {
      it('should close project modal', () => {
        expect(
          ProjectReducer({
            isModalOpen: true
          }, {
            type: newProjectModalActions.CLOSE_NEW_PROJECT,
            isModalOpen: false
          })
        ).toEqual({
          isModalOpen: false
        })
      })
    })

    describe('REQUEST_CREATE_PROJECT', () => {
      it('should close modal and open whole loading', () => {
        expect(
          ProjectReducer({
            projects: ["project1", "project2"],
            isModalOpen: true,
          }, {
            type:newProjectModalActions.REQUEST_CREATE_PROJECT,
            project: "project3"
          })
        ).toEqual({
          projects: ["project1", "project2"],
          isModalOpen: true,
        })
      })
    })

    describe('RECEIVE_CREATE_PROJECT', () => {
      it('should return projects', () => {
        expect(
          ProjectReducer({
            projects: ["project1", "project2"],
          },{
            type: newProjectModalActions.RECEIVE_CREATE_PROJECT,
            project: "project3"
          })
        ).toEqual({
          projects: ["project1", "project2", "project3"],
          isModalOpen: false,
        })
      })
    })
  })
})
