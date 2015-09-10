import React from 'react';
import Router from 'react-router';

// TODO: 初期のオブジェクト渡してもらう方法はなにか考えたほうがいいなぁ
// TODO: react-routerで後でソースや実行関数を分ける
// TODO: ある程度できたらreduxで状態管理する
var App = React.createClass({
    getInitialState: function() {
        return {
            newText: "",
            items: []
        };
    },
    addItem: function() {
        this.setState({
            items: [{text: this.state.newText}].concat(this.state.items),
            newText: ""
        });
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
                    return <div className="item" data-index={index}>
                        {item.text}
                    </div>;
                }, this)}
            </div>
        </div>;
    }

});

React.render(<App name="React" />, document.getElementById('board'));
