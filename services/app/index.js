import { BrowserRouter as Router, Route, Link, Switch } from "react-router-dom";
import ReactDOM from 'react-dom'
import React from 'react'

import { Dashboard, Fontified } from './src/Dashboard'
import { Signup } from './src/SignUp'
import { Login } from './src/Login'
import { Header } from './src/Header'
import { Account } from './src/Account'

import './styles/Main.css'

const HIDDENALPHABET_DOCS_URL = "https://hidden-alphabet.readme.io/docs";

const LandingBackground = props =>
  <div id='landing-background-conntainer'>
    <div id='landing-background'>
      {props.children}
    </div>
  </div>

const CallToActionButton = _ =>
  <div style={{ color: "white", paddingRight: "20px" }}>
    <Link to="/login" style={{ display: 'inline-block' }}>
      <input
      id='call-to-action'
      type="button" 
      style={{ 
        color: "white", 
        width: '200px', 
        height: '40px',
        border: 'none',
        borderRadius: '10px 10px 10px 10px', 
        backgroundColor: '#1aab81',
        boxShadow: '2px',
        fontSize: '13px'
      }} 
      value="Get Started"
      >
      </input>
    </Link>
  </div>

const CallToAction = _ =>
  <div 
  style={{ 
    display: 'flex', 
    flexDirection: 'column', 
    alignItems: 'left', 
    justifyContent: 'center',
    textAlign: 'left',
    paddingTop: '80px',
    paddingLeft: '200px',
    color: 'white'
  }}>
    <p style={{ fontSize: '30px', fontWeight: 'bold', margin: '3' }}>
    Understand your users faster and easier
    </p>
    <h3 style={{ fontSize: '15px', margin: '3' }}>
    Hidden Alphabet is google trends for your business. <br/>
    Know what your users think, today.</h3>    
    <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'center' }}>
      <p style={{ paddingRight: '10px' }}>
      <li>Simple APIs.</li>
      <li>Impactful insights.</li>
      <li>Empower your apps today.</li>
      </p>
    </div>
    <CallToActionButton />
    <Link to="/login">or login.</Link>
  </div>

const Home = _ =>
  <LandingBackground>
    <Fontified>
      <Header>
        {/*}
        <div style={{ color: "white", paddingRight: "20px" }}>
          <Link to="/signup" style={{ color: "white", paddingRight: "10px" }}>Sign Up</Link>{" "}
          <Link to="/login" style={{ color: "white" }}>Login</Link>
        </div>
        */}
      </Header>
      <div style={{ display: 'inline-block' }}>
        <CallToAction />
      </div>      
    </Fontified>
  </LandingBackground>

const App = _ =>
   <Router>
      <Route exact path="/" component={Home} />
      <Route path="/signup" component={Signup} />
      <Route path="/login" component={Login} />
      <Route path="/dashboard" component={Dashboard} />
      <Route path="/account" component={Account} />
    </Router>

ReactDOM.render(<App />, document.getElementById('app'))
