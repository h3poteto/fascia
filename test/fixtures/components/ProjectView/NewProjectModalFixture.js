import expect from 'expect'

export function initState() {
  return {
    isModalOpen: false,
    newProject: {
      title: "title",
      description: ""
    },
    repositories: [{
      id: 1,
      full_name: "repo1"
    }, {
      id: 2,
      full_name: "repo2"
    }],
  }
}
