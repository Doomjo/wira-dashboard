<template>
  <div>
    <PlayerDashboard>

    </PlayerDashboard>
  </div>
</template>

<script>
import PlayerDashboard from './components/PlayerDashboard.vue';
export default {
  name: 'App',
  components: {
    PlayerDashboard
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: bg-white;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
