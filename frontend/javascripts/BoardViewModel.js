import React, { Component } from 'react';
import Request from 'superagent';
import BoardView from './BoardView.jsx';

class BoardViewModel extends Component {
  constructor(props) {
    super(props)
    this.state = {
      isModalOpen: false,
      newProject: "",
      projects: [],
      repositories: [],
      selectedRepository: []
    };
  }
  componentWillMount() {
    var self = this;
    Request
      .get('/projects/')
      .end(function(err, res) {
        if (res.body != null) {
          self.setState({
            newProject: "",
            projects: res.body
          });
        }
      });
  }
  componentDidMount() {
    var self = this;
    Request
      .get('/github/repositories')
      .end(function(err, res) {
        if (res.body != null) {
          self.setState({
            repositories: res.body
          });
        }
      });
  }
  newProject() {
    this.setState({isModalOpen: true});
  }
  closeModal() {
    this.setState({isModalOpen: false});
  }
  createProject() {
    if (this.state.newProject != "") {
      var self = this;
      Request
        .post('/projects/')
        .type('form')
        .send({title: this.props.newProject, repository: this.props.selectedRepository})
        .end(function(err, res) {
          self.setState({
            isModalOpen: false,
            projects: self.state.projects.concat([{Id: res.body.Id, UserId: res.body.UserId, Title: res.body.Title}]),
            newProject: ""
          });
        });
    }
  }
  updateNewText(ev) {
    this.setState({
      newProject: ev.target.value
    });
  }
  changeSelectRepository(ev) {
    this.setState({
      newProject: ev.target.options[ev.target.selectedIndex].text,
      selectedRepository: ev.target.value
    });
  }
  render() {
    return BoardView(this, this.state.newProject, this.state.projects, this.state.repositories, this.state.isModalOpen, this.state.selectedRepository);
  }
}

export default BoardViewModel;
