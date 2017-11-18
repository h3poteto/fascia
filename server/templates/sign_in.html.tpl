{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="session">
  {% include "layouts/_login_header.html.tpl" %}
  <div class="main">
    <article>
      <div class="title">
        <h3>Access your dashboard</h3>
      </div>
      <div class="content">
        <div class="sign-in-board">
          <form action="/sign_in" method="post" role="form" name="sign_in" id="sign_in" class="pure-form pure-form-stacked">
            <fieldset>
              <input name="token" type="hidden" value="{{ token }}" />
              <div class="pure-control-group control-group fascia-form-icon-wrapper">
                <input class="form-control" name="email" type="email" placeholder="email" />
                <div class="fascia-form-icon"><i class="fa fa-user" ></i></div>
              </div>

              <div class="pure-control-group control-group fascia-form-icon-wrapper">
                <input class="form-control" name="password" type="password" placeholder="password" />
                <div class="fascia-form-icon"><i class="fa fa-key" ></i></div>
              </div>

              <div class="pure-controls control-group">
                <a href="/passwords/new">Forgot your password?</a>
              </div>
              <div class="pure-controls control-group">
                <button class="pure-button pure-button-primary session-button" type="submit">SignIn</button>
              </div>
            </fieldset>
          </form>
          <a href={{ publicURL }}><span class="pure-button button-success session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github</span></a>
          <span class="message">If you want to management private repository, <a href={{ privateURL }}>please click here</a>.</span>
          <p><a href="/sign_up"><span class="pure-button button-secondary session-button">SignUp</span></a></p>
        </div>
      </div>
    </article>
  </div>
</div>
{% endblock %}
