import React from 'react'

const h = React.createElement

class MemberElement extends React.Component {
//    propTypes: {
//        name: React.PropTypes.string.isRequired,
//        netid: React.PropTypes.string.isRequired,
//        },

  removeMember (event) {
    fetch('/api/group/physics-grp-test/remove/' + this.props.NetId,
     {method: 'POST', credentials: 'same-origin'}).then()
//    this.props.onGroupModified()
  }

  render () {
    return (
      h('li', {className: 'member'},
      h('a', {
        className: 'member-rm-btn',
        onClick: e => this.removeMember(e)
      }, ' × '),
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

export default class App extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      groupname: '',
      members: []
    }
  }

  fetchState () {
    fetch('/api/group/physics-grp-test/members', {credentials: 'same-origin'})
      .then(r => r.json())
      .then(data => {
        this.setState({groupname: data.Name, members: data.Users}, () => this.render())
      })
  }

  componentDidMount () {
    this.fetchState()
  }

  render () {
    return (
      h('div', null,
        h('h1', null, this.state.groupname),
        h('ul', {className: "member-list"}, this.state.members.map(member => h(MemberElement, member))),
        h(AddMemberElement, {onGroupModified: () => this.fetchState()}))
    )
  }
}
