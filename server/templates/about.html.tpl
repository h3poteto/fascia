{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="about">
  {% include "layouts/_login_header.html.tpl" %}
  <div class="main-board">
    <h1>Simple Task Management</h1>
    <h2>Fascia is a free task management solution.</h2>
    <article>
      <div class="title">
        <h3>Sign Up - It's Free</h3>
      </div>
      <div class="content">
        <div class="sign-up-board">
          <a href="/oauth/sign_in"><span class="pure-button button-success session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github</span></a>
        </div>
      </div>
    </article>
  </div>
  <div class="main-area">
    <div class="store-link">
      <a href="https://itunes.apple.com/us/app/fascia/id1197999536?mt=8" style="display:inline-block;overflow:hidden;background:url(https://linkmaker.itunes.apple.com/en-us/badge-lrg.svg?releaseDate=2018-07-28&kind=iossoftware&bubble=ios_apps) no-repeat;width:135px;height:40px;"></a>
    </div>
    <div class="github">
      <h2>Sync GitHub</h2>
      <div class="content">
        <p>Fascia can sync your GitHub repositories, no matter if private repository.</p>
        <p>Fascia's tasks are reflected in GitHub issues, and get Fascia's tasks from Github issues.</p>
        <table class="github-images">
          <tr>
            <td>
              <img src="/lp/images/fascia-ss-2.png" class="fascia-ss-2">
            </td>
            <td>
              <img src="/lp/images/github-ss-1.png" class="github-ss-1">
            </td>
          </tr>
        </table>

        <p>In addition, you can create tasks which are not related to GitHub labels.</p>
        <p>At first, you create tasks which do not belong to any GitHub labels, and then you can move task(GitHub issue) to a GitHub label.</p>
      </div>
    </div>
    <div class="not-github">
      <h2>Projects are not synced GitHub</h2>
      <div class="content">
        <p>Fascia can manage projects are not related GitHub repositories.</p>
        <p>Therefore, please manage your private tasks, for example, shopping list and domestic work.</p>
        <img src="/lp/images/fascia-ss-3.png" class="fascia-ss-3">
      </div>
    </div>
  </div>
  <div class="footer">
    <p>&copy; Copyright 2016, <a href="https://twitter.com/h3_poteto">@h3_poteto</a></p>
  </div>
</div>
{% endblock %}
