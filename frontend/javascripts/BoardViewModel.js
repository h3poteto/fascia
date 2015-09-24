import React from 'react';
import Request from 'superagent';
import BoardView from './BoardView.jsx';

var BoardViewModel = React.createClass({
  getInitialState: function() {
    return {
      isModalOpen: false,
      newText: "",
      items: []
    };
  },
  componentDidMount: function() {
    var self = this;
    Request
      .get('/projects/')
      .end(function(err, res) {
        if (self.isMounted() && res.body != null) {
          self.setState({
            newText: "",
            items: res.body
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
    if (this.state.newText != "") {
      var self = this;
      Request
        .post('/projects/')
        .type('form')
        .send({title: this.state.newText})
        .end(function(err, res) {
          self.setState({
            isModalOpen: false,
            items: self.state.items.concat([{Id: res.body.Id, UserId: res.body.UserId, Title: res.body.Title}]),
            newText: ""
          });
        });
    }
  },
  updateNewText: function(ev) {
    this.setState({
      newText: ev.target.value
    });
  },

  render: function() {
    return BoardView(this, this.state.newText, this.state.items, this.state.isModalOpen);
  }
});

export default BoardViewModel;
