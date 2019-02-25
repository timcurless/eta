class App extends React.Component {
  render() {
    return <Home />;
  }
}

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      health: [],
      database: [],
      users: []
    };
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleLoadUsers = this.handleLoadUsers.bind(this);
  }

  serverRequest(path) {
    $.get("http://localhost:3000" + path, res => {
      this.setState({
        health: res.health,
        database: res.database
      });
    });
  }

  componentDidMount() {
    this.serverRequest("/api/health");
    this.handleLoadUsers();
  }

  handleSubmit(data) {
    $.post("http://localhost:3000/api/users", JSON.stringify(data), res => {
      this.handleLoadUsers();
    })
  }

  handleLoadUsers() {
    $.get("http://localhost:3000/api/users", res => {
      this.setState({
        users: res
      });
    });
  }

  render() {
    return (
      <div className="container">
        <div className="jumbotron text-center">
          <h1>{'Welcome to \u03B7!'}</h1>
          <hr />
          <div className="container text-left">
            <div className="row">
              <div className="col-md">
                <h3>Vault Status</h3>
                <div className="table-responsive">
                    <pre>{ JSON.stringify(this.state.health, null, 2) }</pre>
                </div>
              </div>
              <div className="col-md">
                <h3>Database Status</h3>
                  <pre>{ JSON.stringify(this.state.database, null, 2) }</pre>
              </div>
            </div>
          </div>
          <hr />

          <div className="row justify-content-center">
            <UserForm onSubmit={this.handleSubmit}/>
          </div>
          <hr />
          <div className="row justify-content-center">
            <h3>User List</h3>
            <UserTable users={this.state.users}/>
          </div>
        </div>
      </div>
    );
  }
}

class UserForm extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      first_name: '',
      last_name: '',
      email: '',
      password: ''
    };

    this.handleChange = this.handleChange.bind(this);
    this.submitNewUser = this.submitNewUser.bind(this);
  }

  submitNewUser(e) {
    const { onSubmit } = this.props;
    e.preventDefault();
    onSubmit(this.state);
  }

  handleChange(e) {
    this.setState({
      [e.target.name]: e.target.value
    });
  }

  render() {
    return (
      <form style={{width: '40em'}}
            onSubmit={this.submitNewUser}>
        <h3>Create New User</h3>
        <div className="form-group row">
          <input type="text"
                 className="form-control"
                 name="first_name"
                 value={this.state.first_name}
                 onChange={this.handleChange}
                 placeholder="Enter First Name" />
        </div>
        <div className="form-group row">
          <input type="text"
                 className="form-control"
                 name="last_name"
                 value={this.state.last_name}
                 onChange={this.handleChange}
                 placeholder="Enter Last Name" />
        </div>
        <div className="form-group row">
          <input type="email"
                 className="form-control"
                 name="email"
                 value={this.state.email}
                 onChange={this.handleChange}
                 placeholder="Enter Email Address" />
        </div>
        <div className="form-group row">
          <input type="password"
                 className="form-control"
                 name="password"
                 value={this.state.password}
                 onChange={this.handleChange}
                 placeholder="Enter Password" />
          <small id="passwordHelpBlock" class="form-text text-muted">
            Your password must be 8-20 characters long, contain letters and numbers, and must not contain spaces, special characters, or emoji.
          </small>
        </div>
        <button type="submit"
                className="btn btn-primary">Create User</button>
      </form>
    );
  }
}

class UserTable extends React.Component {
  render() {
    return (
      <table className="table table-sm">
        <thead className="thead-dark">
          <tr>
            <th>First Name</th>
            <th>Last Name</th>
            <th>Email Address</th>
            <th>AheadAviation Rewards ID</th>
          </tr>
        </thead>
        <tbody>
          {
            this.props.users && this.props.users.map((user) => {
              return (
                <tr id={user.ID}>
                  <td>{user.first_name}</td>
                  <td>{user.last_name}</td>
                  <td>{user.email}</td>
                  <td>{user.rewards_id}</td>
                </tr>
              );
            })
          }
        </tbody>
      </table>
    );
  }
}

ReactDOM.render(<App />, document.getElementById("app"));
