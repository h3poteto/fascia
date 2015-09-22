import React from 'react';
import {Router, Route, Link} from 'react-router';
import Request from 'superagent';


// TODO: ある程度できたらreduxで状態管理する
var Board = React.createClass({
    getInitialState: function() {
        return {
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
    addItem: function() {
        if (this.state.newText != "") {
            var self = this;
            Request
                .post('/projects/')
                .type('form')
                .send({title: this.state.newText})
                .end(function(err, res) {
                    self.setState({
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
        return <div>
            <input type="text" value={this.state.newText} onChange={this.updateNewText} placeholder="Project Name" className="form-control" />
            <button onClick={this.addItem} className="pure-button pure-button-primary fascia-button" type="button">CreateProject</button>
            <div className="items">
                {this.state.items.map(function(item, index) {
                    return <div className="item" data-id={item.Id}>
                        {item.Title}
                    </div>;
                }, this)}
            </div>
        </div>;
    }

});


var routes = (
  <Router>
    <Route path="/" component={Board}>
    </Route>
  </Router>
);

React.render(routes, document.getElementById("board"));


var Menu = React.createClass({
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
        return <div className="top-nav">
            <div className="pure-menu pure-menu-horizontal">
            <span className="pure-menu-heading">fascia</span>
            <ul className="pure-menu-list fascia-menu-list">
            <li className="pure-menu-item fascia-menu-item"><a href="/" className="pure-menu-link">projects</a></li>
            </ul>
            <ul className="pure-menu-list fascia-menu-list right-align-list">
            <li className="pure-menu-item fascia-menu-item"><a href="#" className="pure-menu-link">settings</a></li>
            </ul>
            </div>
            </div>;
        }
});

var menuRoutes = (
    <Router>
      <Route path="/" component={Menu}>
      </Route>
    </Router>
);

React.render(menuRoutes, document.getElementById("top"));
