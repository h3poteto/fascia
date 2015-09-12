import React from 'react';
import Router from 'react-router';
import Request from 'superagent';

// TODO: react-routerで後でソースや実行関数を分ける
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
            <button onClick={this.addItem} className="pure-button pure-button-primary" type="button">CreateProject</button>
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

React.render(<Board name="React" />, document.getElementById('board'));
