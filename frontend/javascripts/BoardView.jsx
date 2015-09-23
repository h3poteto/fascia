import React from 'react';

var BoardView = function(component, text, projects) {
  return <div id="board">
    <input type="text" value={text} onChange={component.updateNewText} placeholder="Project Name" className="form-control" />
    <button onClick={component.addItem} className="pure-button pure-button-primary fascia-button" type="button">CreateProject</button>
    <div className="items">
      {projects.map(function(item, index) {
        return <div className="item" data-id={item.Id}>
        {item.Title}
        </div>;
       }, component)}
    </div>
  </div>;

 }

export default BoardView;
