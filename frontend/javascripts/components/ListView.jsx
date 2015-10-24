import React from 'react';
import Modal from 'react-modal';

const customStyles = {
  overlay : {
    position          : 'fixed',
    top               : 0,
    left              : 0,
    right             : 0,
    bottom            : 0,
    backgroundColor   : 'rgba(255, 255, 255, 0.5)'
  },
  content : {
    position : 'fixed',
    top : '50%',
    left : '50%',
    right : 'auto',
    bottom : 'auto',
    marginRight : '-50%',
    transform : 'translate(-50%, -50%)'
  }
};

export default class ListView extends React.Component {
  constructor(props) {
    super(props);
  }

  componentWillMount() {
    this.props.fetchLists(this.props.params.projectId);
  }

  render() {
    console.log(this.props);
    const { isModalOpen, newList, lists } = this.props.ListReducer
    return (
      <div id="lists">
        <Modal
          isOpen={isModalOpen}
          onRequestClose={this.props.closeNewListModal}
          style={customStyles}
        >
          <div className="list-form">
            <form className="pure-form pure-form-stacked">
              <fieldset>
                <legend>Create List</legend>
                <label htmlFor="title">Title</label>
                <input id="title" name="title" type="text" value={newList.title} onChange={this.props.updateNewListTitle} placeholder="List Name" className="form-control" />
                <button onClick={e => this.props.fetchCreateList(this.props.params.projectId, newList.title)} className="pure-button pure-button-primary" type="button">Create List</button>
              </fieldset>
            </form>
          </div>
        </Modal>
        <div className="items">
          {lists.map(function(item, index) {
            return (
              <div className="fascia-card pure-button button-secondary" data-id={item.Id}>
              <span className="card-title">{item.Title}</span>
              </div>
            );
           }, this)}
              <button onClick={this.props.openNewListModal} className="pure-button button-large fascia-new-list" type="button">New</button>
        </div>
      </div>
    );
  }
}
