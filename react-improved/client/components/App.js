import React from 'react'
import { withNotie } from 'react-notie'

const h = React.createElement

function sleep (ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

function stripPhysicsPrefix(id) {
  return id.replace('physics-grp-', '')
}

class MemberElement extends React.Component {
  constructor (props) {
    super(props)
  }
//    propTypes: {
//        name: React.PropTypes.string.isRequired,
//        netid: React.PropTypes.string.isRequired,
//        },

  removeMember (event) {
    fetch('/api/group/'+ this.props.selectedgroup + '/remove/' + this.props.NetId,
     {method: 'POST', credentials: 'same-origin'}).then((r) => r.json()).then((msg) => {
       if (msg.Result == 'Success') {
         this.props.notie.success(msg.Message)
       }
       else if (msg.Result == 'Warn') {
         this.props.notie.warn(msg.Message)
       }
       else {
         this.props.notie.error(msg.Message)
       }
     })
    this.props.onGroupModified()
  }

  render () {
    return (
      h('li', {className: 'member'},
      h('a', {
        className: 'member-rm-btn',
        onClick: e => this.removeMember(e)
      }, '  ✕  '),
      h('span', {className: 'member-name'}, this.props.Name),
      h('span', {className: 'member-netid'}, '  (' + this.props.NetId + ')')
      )
    )
  }
}

class AddMemberElement extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      netid: '' // holds the netid as the user is typing it
    }
  }

  handleSubmit (event) {
    event.preventDefault()
    fetch('/api/group/' + this.props.selectedgroup + '/add/' + this.state.netid,
     {method: 'POST', credentials: 'same-origin'}).then((r) => r.json()).then((msg) => {
       if (msg.Result == 'Success') {
         this.props.notie.success(msg.Message)
       }
       else if (msg.Result == 'Warn') {
         this.props.notie.warn(msg.Message)
       }
       else {
         this.props.notie.error(msg.Message)
       }
     })
    this.setState({netid: ''})
    this.props.onGroupModified()
  }

  render () {
    return (
      h('form', {
        className: 'add-member-form',
        onSubmit: event => this.handleSubmit(event),
        onChange: event => this.setState({netid: event.target.value})
      },
        h('input', {
          type: 'text',
          placeholder: 'NetId (required)',
          value: this.state.netid
        }),
        h('button', {
          type: 'submit'
        }, 'Add')
      )
    )
  }
}

class GroupViewComponent extends React.Component {
  constructor (props) {
    super(props)
  }

  render() {
    return(
        h('div', {className: 'managed-groups-container'}, 
          h('ul', {className: 'group-list member-list'}, 
            this.props.groups.map(group=>h('li', {onClick: () => this.props.onGroupSelected(group)}, stripPhysicsPrefix(group))))
        )
    )
  }
}

class MemberViewComponent extends React.Component {
  constructor (props) {
    super(props)
  }

  render() {
    if (this.props.selectedgroup == '') { // i.e., nothing selected
      return (
        h('div', {id: 'group-members-placeholder'},
          h('img', {src: "group_not_selected.svg"}),
          h('h2', {className: "sans-font"}, "No group selected"),
          h('h3', {className: "sans-font"}, "Choose a group from the left to get started!")
        )
      )
    }
    else {
      return(
          h('div', {id: 'group-members-container'},
            h('h1', {className: 'sans-font'}, stripPhysicsPrefix(this.props.selectedgroup)),
            h('ul', {className: 'member-list'}, this.props.members.map(member => {
              Object.assign(member, {
                selectedgroup: this.props.selectedgroup,
                onGroupModified: this.props.onGroupModified,
                notie: this.props.notie
              })
              return h(MemberElement, member)
            })),
            h(AddMemberElement, {
              selectedgroup: this.props.selectedgroup,
              onGroupModified: this.props.onGroupModified,
              notie: this.props.notie
            }))
      )
    }
  }
}

class App extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      groups: [],
      selectedgroup: '',
      members: []
    }

    this.fetchState = this.fetchState.bind(this)
    this.selectGroup = this.selectGroup.bind(this)
  }

  selectGroup(group) {
    this.state.selectedgroup = group
    this.fetchState()
  }

  fetchState () {
    fetch('/api/user/managed', {credentials: 'same-origin'})
      .then(r => r.json())
      .then(data => {
        this.setState({groups: data})
    })

    if (this.state.selectedgroup != '') {
      sleep(500).then(() => {
        fetch('/api/group/'+ this.state.selectedgroup + '/members', {credentials: 'same-origin'})
          .then(r => r.json())
          .then(data => {
            this.setState({selectedgroup: data.Name, members: data.Users})
          })
      })
    }
  }

  componentDidMount () {
    this.fetchState()
  }

  render () {
    return (
      h('div', {id: 'main-app'}, 
        h(GroupViewComponent, {
          groups: this.state.groups,
          onGroupSelected: this.selectGroup
        }),
        h(MemberViewComponent, {
          members: this.state.members,
          selectedgroup: this.state.selectedgroup,
          onGroupModified: this.fetchState,
          notie: this.props.notie
        })
      )
    )
  }
}

export default withNotie(App)