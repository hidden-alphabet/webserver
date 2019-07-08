import { BrowserRouter as Router, Route, Link, Switch } from "react-router-dom";
import ReactDOM from 'react-dom'
import React from 'react'

import { Dashboard, Fontified } from './src/Dashboard'
import { Signup } from './src/SignUp'
import { Login } from './src/Login'
import { Header } from './src/Header'
import { Account } from './src/Account'

import './styles/Main.css'

import comparison from './static/comparison.png'
import infrastructure from './static/infrastructure.png'

const LandingBackground = props =>
  <div id='landing-background-container'>
    <div id='landing-background'>
      {props.children}
    </div>
  </div>

const CallToActionButton = _ =>
  <div style={{ color: "white", paddingRight: "20px" }}>
    <a href="https://hidden-alphabet.readme.io/docs" style={{ display: 'inline-block' }}>
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
    </a>
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
    Know what your users think of you, today.</h3>    
    <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'center' }}>
      <p style={{ paddingRight: '10px' }}>
      <li>Simple APIs.</li>
      <li>Impactful insights.</li>
      <li>Empower your apps today.</li>
      </p>
    </div>
    <CallToActionButton />
    {/*<Link to="/login">or login.</Link> */}
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
      
      <div style={{ display: 'flex', flexDirection: 'column', marginTop: '50px' }}>
        <div style={{ 
          height: '500px', 
          width: '100%', 
          backgroundColor: 'white', 
          display: 'flex', 
          flexDirection: 'row', 
          justifyContent: 'center',
          alignItems: 'center'
        }}>
          <img style={{ backgroundColor: 'white' }} src={comparison}></img>
        </div>
        <div style={{ 
          height: '500px', 
          width: '100%', 
          backgroundColor: '#ac38fb', 
          display: 'flex', 
          flexDirection: 'row', 
          justifyContent: 'space-between',
          alignItems: 'center'
        }}>
          <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', alignItems: 'items', margin: '30px' }}>
            <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', backgroundColor: 'white', margin: '70px', padding: '60px 40px 60px 40px', borderRadius: '6%' }}>
              <h3>Donaldo Celaj</h3>
              <p>NLP Developer</p>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', backgroundColor: 'white',  margin: '70px', padding: '60px 40px 60px 40px', borderRadius: '6%' }}>
              <h3>Cole Hudson</h3>
              <p>ML Backend architect</p>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', backgroundColor: 'white',  margin: '70px', padding: '60px 40px 60px 40px', borderRadius: '6%' }}>
              <h3>Mac Scheffer</h3>
              <p>NLP Developer</p>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', backgroundColor: 'white',  margin: '70px', padding: '60px 40px 60px 40px', borderRadius: '6%'}}>
              <h3>Chris Louie</h3>
              <p>ML Backend Archiect</p>
            </div>
          </div>
        </div>
        <div style={{ 
          height: '800px', 
          width: '100%', 
          backgroundColor: 'white', 
          display: 'flex', 
          flexDirection: 'row', 
          justifyContent: 'center',
          alignItems: 'center'
        }}>
          <h1>AWS Archiecture</h1> 
          <img style={{ wifth: 'auto', heigth: 'auto', maxWidth: '900px', maxHeight: '1000px', backgroundColor: 'white' }} src={infrastructure}></img>
        </div>
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
