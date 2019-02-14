class App extends React.Component {
  render() {
    return <Home />;
  }
}

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      health: []
    };
  }
  serverRequest() {
    $.get("http://localhost:3000/api/health", res => {
      this.setState({
        health: res
      });
    });
  }
  componentDidMount() {
    this.serverRequest();
  }
  render() {
    return (
      <div className="container">
        <div className="row">
          <div className="jumbotron text-center">
            <h1>Welcome to Eta!</h1>
            <div className="row">
              <h4>Vault Status:</h4>
              <div className="container">
                {JSON.stringify(this.state.health)}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById("app"));
