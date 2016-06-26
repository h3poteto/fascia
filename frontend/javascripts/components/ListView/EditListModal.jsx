import React from 'react'
import Modal from 'react-modal'

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
}

export default class EditListModal extends React.Component {
  constructor(props) {
    super(props)
  }

  listAction(project, listOptions, selectedList, selectedListOption) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return null
    } else {
      return (
        <div>
          <label htmlFor="action">action</label>
          <select id="action" name="action" type="text" onChange={this.props.changeSelectedListOption} className="form-control" value={selectedListOption ? selectedListOption.ID : (selectedList ? selectedList.ListOptionID : 0)}>
            <option value="0">nothing</option>
            {listOptions.map(function(option, index) {
               return <option key={index} value={option.ID}>{option.Action}</option>
             }, this)}
          </select>
        </div>
      )
    }
  }

  render() {
    return (
      <Modal
          isOpen={this.props.isListEditModalOpen}
          onRequestClose={this.props.closeEditListModal}
          style={customStyles}
      >
        <div className="list-form">
          <form className="pure-form pure-form-stacked">
            <fieldset>
              <legend>Edit List</legend>
              <label htmlFor="title">Title</label>
              <input id="title" name="title" type="text" value={this.props.selectedList !=null ? this.props.selectedList.Title : ''} onChange={this.props.updateSelectedListTitle} className="form-control" />
              <label htmlFor="color">Color</label>
              <input id="color" name="color" type="text" value={this.props.selectedList !=null ? this.props.selectedList.Color : ''} onChange={this.props.updateSelectedListColor} className="form-control" />
              {this.listAction(this.props.project, this.props.listOptions, this.props.selectedList, this.props.selectedListOption)}
              <div className="form-action">
                <button onClick={e => this.props.fetchUpdateList(this.props.projectID, this.props.selectedList, this.props.selectedListOption)} className="pure-button pure-button-primary" type="button">Update List</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}
