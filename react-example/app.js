var h = React.createElement;

var MemberElement = React.createClass({
    propTypes: {
        name: React.PropTypes.string.isRequired,
        netid: React.PropTypes.string.isRequired,
        },
    render: function() {
        return (
            h('li', {className: 'member'},
            h('span' , {className: 'member-name'}, this.props.Name),
            h('span', {className: 'member-netid'}, '  (' + this.props.NetId + ')')
            )
        )
    },
})

var newMember = {NetId: ""}

var AddMemberElement = React.createClass({
    render : function() {
        return (
            h('form', {
                className: 'add-member-form',
                onChange: function(syntheticEvent) {
                    console.log(syntheticEvent.target.value);
                },
            },
            h('input', {
                type: 'text',
                placeholder: 'NetId (required)',
                value: this.props.NetId,
            }),
            h('button', {type: 'submit'}, 'Add')
            )
        )
    },
})

var groupData = {};

$.ajax({
    async:false,
    url: '/api/group/physics-grp-test/members',
    success: function(data) {
        groupData = data;
        }
})

var MemberItemElements = groupData.Users
.map(function(member) {return h(MemberElement, member)});



var rootElement = h('div', null,
                    h('h1', null, groupData.Name),
                    h('ul', null, MemberItemElements),
                    h(AddMemberElement, newMember)
);

ReactDOM.render(rootElement, document.getElementById("react-app"));