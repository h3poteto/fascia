{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="session">
  {% include "layouts/_login_header.html.tpl" %}
  <div class="main">
    <article>
      <div class="title">
        <h3>Create Acount</h3>
      </div>
      <div class="content">
        <div class="sign-up-board">
          <form action="/sign_up" method="post" role="form" name="sign_up" id="sign_up" class="pure-form pure-form-stacked">
            <fieldset>
              <input name="token" type="hidden" value="{{ token }}" />
              <div class="pure-control-group control-group fascia-form-icon-wrapper">
                <input class="form-control" name="email" type="email" placeholder="email" />
                <div class="fascia-form-icon"><i class="fa fa-user"></i></div>
              </div>

              <div class="pure-control-group control-group fascia-form-icon-wrapper">
                <input class="form-control" name="password" type="password" placeholder="password" />
                <div class="fascia-form-icon"><i class="fa fa-key" ></i></div>
              </div>

              <div class="pure-control-group control-group fascia-form-icon-wrapper">
                <input class="form-control" name="password_confirm" type="password" placeholder="password" />
                <div class="fascia-form-icon"><i class="fa fa-key" ></i></div>
              </div>

              <span class="button-small">
                You have read and agree to the <a href="/privacy_policy">privacy policy</a>
              </span>
              <div class="pure-controls control-group">
                <button class="pure-button pure-button-primary session-button" type="submit">SignUp</button>
              </div>
            </fieldset>
          </form>
          <a href="/oauth/sign_in"><span class="pure-button button-success session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github</span></a>
        </div>
      </div>
    </article>
  </div>
</div>
{% endblock %}
