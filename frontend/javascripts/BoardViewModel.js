import React from 'react';
import Request from 'superagent';
import BoardView from './BoardView.jsx';

var BoardViewModel = React.createClass({
  getInitialState: function() {
    return {
      isModalOpen: false,
      newProject: "",
      projects: [],
      repositories: [],
      selectedRepository: []
    };
  },
  componentWillMount: function() {
    var self = this;
    Request
      .get('/projects/')
      .end(function(err, res) {
        if (self.isMounted() && res.body != null) {
          self.setState({
            newProject: "",
            projects: res.body
          });
        }
      });
  },
  componentDidMount: function() {
    var self = this;
    Request
      .get('/github/repositories')
      .end(function(err, res) {
        if (self.isMounted() && res.body != null) {
          self.setState({
            repositories: res.body
          });
        }
      });
  },
  newProject: function() {
    this.setState({isModalOpen: true});
  },
  closeModal: function() {
    this.setState({isModalOpen: false});
  },
  createProject: function() {
    if (this.state.newProject != "") {
      var self = this;
      Request
        .post('/projects/')
        .type('form')
        .send({title: this.state.newProject})
        .end(function(err, res) {
          self.setState({
            isModalOpen: false,
            projects: self.state.projects.concat([{Id: res.body.Id, UserId: res.body.UserId, Title: res.body.Title}]),
            newProject: ""
          });
        });
    }
  },
  updateNewText: function(ev) {
    this.setState({
      newProject: ev.target.value
    });
  },
  changeSelectRepository: function(ev) {
    this.setState({
      newProject: ev.target.options[ev.target.selectedIndex].text,
      selectedRepository: ev.target.value
    });
  },

  render: function() {
    return BoardView(this, this.state.newProject, this.state.projects, this.state.repositories, this.state.isModalOpen, this.state.selectedRepository);
  }
});

export default BoardViewModel;
