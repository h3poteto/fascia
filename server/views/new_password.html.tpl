{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="session">
  {% include "layouts/_login_header.html.tpl" %}
  <div class="main">
    <article>
      <div class="title">
        <h3>Reset your password</h3>
      </div>
      <div class="content">
        <div class="sign-in-board">
          <form action="/passwords/create" method="post" role="form" name="create" id="create" class="pure-form pure-form-stacked">
            <fieldset>
              <input name="token" type="hidden" value="{{ token }}" />
              <div class="pure-control-group control-group fascia-form-icon-wrapper">
                <input class="form-control" name="email" type="email" placeholder="email" />
                <div class="fascia-form-icon"><i class="fa fa-user" ></i></div>
              </div>
              <div class="pure-controls control-group">
                <button class="pure-button pure-button-primary session-button" type="submit">Reset Password</button>
              </div>
            </fieldset>
          </form>
        </div>
      </div>
    </article>
  </div>
</div>
{% endblock %}
