import React from 'react';
import $ from 'jquery'

const h = React.createElement;


class MemberElement extends React.Component {
//    propTypes: {
//        name: React.PropTypes.string.isRequired,
//        netid: React.PropTypes.string.isRequired,
//        },
    render() {
        return (
            h('li', {className: 'member'},
            h('span' , {className: 'member-name'}, this.props.Name),
            h('span', {className: 'member-netid'}, '  (' + this.props.NetId + ')')
            )
        )
    }
}


class AddMemberElement extends React.Component {

    render() {
        return (
            h('form', {
                className: 'add-member-form',
                onChange: syntheticEvent => null,
                },
              h('input', {
                  type: 'text',
                  placeholder: 'NetId (required)',
              }),
              h('button', {type: 'submit'}, 'Add')
            )
        )
    }
}

var newMember = {NetId: ""}

export default class App extends React.Component {

    constructor(props) {
      super(props);
      this.state = {
        groupname: '',
        members: []
      };
    }

  componentDidMount() {
      fetch('/api/group/physics-grp-test/members', {credentials: 'same-origin'})
        .then(r => r.json())
        .then(data => {
          this.setState({groupname: data.Name});
          this.setState({members: data.Users})
      });
  }

  render() {
    return (
        h('div', null,
          h('h1', null, this.state.groupname),
          h('ul', null, this.state.members.map(member => h(MemberElement, member))),
          h(AddMemberElement, newMember))
    );
  }
}
