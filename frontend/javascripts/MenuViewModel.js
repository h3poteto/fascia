import React from 'react';
import MenuView from './MenuView.jsx';

var MenuViewModel = React.createClass({
  getInitialState: function() {
    return {
      selected: "projects"
    };
  },
  selectMenu: function(selectedClass) {
    this.setState({
      selected: selectedClass
    });
  },
  render: function() {
    return MenuView(this.state.selected);
  }
});

export default MenuViewModel;
