import React from 'react'

const h = React.createElement

function sleep (ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
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
    fetch('/api/group/physics-grp-test/remove/' + this.props.NetId,
     {method: 'POST', credentials: 'same-origin'}).then()
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
      netid: ''
    }
  }

  handleSubmit (event) {
    event.preventDefault()
    fetch('/api/group/physics-grp-test/add/' + this.state.netid,
     {method: 'POST', credentials: 'same-origin'}).then()
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

class MemberViewComponent extends React.Component {
  constructor (props) {
    super(props)
  }

  render() {
    if (this.props.selectedgroup == '') { // i.e., nothing selected
      return (
        h('div', {id: 'group-members-placeholder'})
      )
    }
    else {
      return(
          h('div', {id: 'group-members-container'},
            h('h1', null, this.props.selectedgroup),
            h('ul', {className: 'member-list'}, this.props.members.map(member => {
              Object.assign(member, {onGroupModified: this.props.onGroupModified})
              return h(MemberElement, member)
            })),
            h(AddMemberElement, {onGroupModified: this.props.onGroupModified}))
      )
    }
  }
}

export default class App extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      groups: [],
      selectedgroup: 'physics-grp-test',
      members: []
    }

    this.fetchState = this.fetchState.bind(this)
  }

  fetchState () {
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
        h('div', {className: 'managed-groups-container'}, 
          h('ul', {className: 'group-list'})

        ),
        h(MemberViewComponent, {
           members: this.state.members,
           selectedgroup: this.state.selectedgroup,
           onGroupModified: this.fetchState
          })
      )
    )
  }
}
