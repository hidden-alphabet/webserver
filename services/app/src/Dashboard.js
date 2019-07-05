import { BrowserRouter as Router, Route, Link, Redirect } from 'react-router-dom' 
import React from 'react';

import { Header } from './Header'

const DMSansFont = _ =>
  <link href="https://fonts.googleapis.com/css?family=DM+Sans&display=swap" rel="stylesheet">
  </link>

const Fontified = props => 
  <div style={{ 'fontFamily': 'DM Sans' }}>
    <DMSansFont />    
    {props.children}
  </div>

const HeaderLinks = _ =>
  <div style={{ display: "flex", flexDirection: "row", justifyContent: "flex-end" }}>
    <Link to="/account" style={{ paddingRight: "10px" }}>Account</Link>
    <Link to="/">Logout</Link>
  </div>

const AnalyticsTab = _ =>
  <>
    <h2>Analytics</h2>
    <hr/>
    <p>
      Our technical wizards's are off gathering data just for you,
      come back later to see what they've found! 
      <br/>
      Or,{" "} <a href="/events/subscribe">be notified by email</a> 
      {" "}when it's ready!
    </p>
  </>

class EventsTab extends React.Component {
  constructor(props) {
    super(props)

    // let res = axios.post('/user/auth/test')
    /*
      if(!res.isAuthenticated) {
        return
      }
    */

    this.eventSource = new EventSource("events");

    this.state = {}
    this.state.data = []

    this.columns = [
      {
        Header: "Resource",
        accessor: "resource"
      },
      {
        Header: "Date",
        accessor: "date"
      },
      {
        Header: "Change (%)",
        accessor: "change"
      },
      {
        Header: "Direction",
        accessor: "direction"
      }
    ];
  }

  updateFlightState(flightState) {
    let newData = this.state.data.map(item => {
      if (item.flight === flightState.flight) {
        item.state = flightState.state;
      }
      return item;
    });
   
    this.setState(Object.assign({}, { data: newData }));
  }

  componentDidMount() {
    this.eventSource.onmessage = e => this.updateFlightState(JSON.parse(e.data));
  }

  render() {
    return (
        <div style={{ margin: '0' }}>
        <h2>Events</h2>
        <hr/>
        <p>
            Nothing to see here
        </p>
      </div>
    )
  }
}

class Dashboard extends React.Component {
  render() {
    return (
      <Fontified>
        <Header>
          <HeaderLinks />
        </Header>
        <div id='container' style={{ display: 'flex', flexDirection: 'row', justifyConten: 'space-between', alignItems: 'flex-start', height: '100%', padding: '0', margin: '0' }}>
         <Router>
            <div id='sidebar' style={{ display: 'flex', flexDirection: 'column', paddingRight: '50px', paddingLeft: '10px', backgroundColor: 'lightgrey', height: '100%' }}>
              <Link to="/dashboard/sentiment">Sentiment</Link>
              <Link to="/dashboard/events">Events</Link>
            </div>
            <div id='display' style={{ display: 'flex', flexDirection: 'column', height: '100%', paddingLeft: '20px' }}>
              <Redirect from="/dashboard" exact to="/dashboard/sentiment" />
              <Route exact path="/dashboard/sentiment" component={AnalyticsTab} />
              <Route path="/dashboard/events" component={EventsTab} />
            </div>
          </Router>
        </div>
      </Fontified>
    )
  }
}

export {
  Dashboard,
  Fontified
}