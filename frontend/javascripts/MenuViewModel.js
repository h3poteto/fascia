import React, { Component } from 'react';
import MenuView from './MenuView.jsx';

class MenuViewModel extends Component {
  constructor(props) {
    super(props)
    this.state = {
      selected: "projects"
    };
  }
  selectMenu(selectedClass) {
    this.setState({
      selected: selectedClass
    });
  }
  render() {
    return MenuView(this.props, this.state.selected);
  }
}

export default MenuViewModel;
