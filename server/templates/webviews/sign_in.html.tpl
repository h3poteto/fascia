{% extends "../layouts/webview.html.tpl" %}

{% block content %}
<div class="session">
  <div class="main">
    <div class="sign-in-board">
      <form action="/webviews/sign_in" method="post" role="form" name="sign_in" class="pure-form pure-form-stacked">
        <fieldset>
          <input name="token" type="hidden" value="{{ token }}" />
          <div class="pure-control-group control-group fascia-form-icon-wrapper">
            <input class="form-control" name="email" type="email" placeholder="email" />
            <div class="fascia-form-icon"><i class="fa fa-user"></i></div>
          </div>

          <div class="pure-control-group control-group fascia-form-icon-wrapper">
            <input class="form-control" name="password" type="password" placeholder="password" />
            <div class="fascia-form-icon"><i class="fa fa-key"></i></div>
          </div>

          <div class="pure-controls control-group">
            <button class="pure-button pure-button-primary session-button" type="submit">SignIn</button>
          </div>
        </fieldset>
      </form>
      <a href={{ publicURL }}><span class="pure-button button-success session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github</span></a>
      <span class="message">If you want to management private repository, <a href={{ privateURL }}>please click here</a>.</span>
    </div>
  </div>
</div>
{% endblock %}
