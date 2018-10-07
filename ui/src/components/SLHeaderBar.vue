<template>
  <div class="header">
    <div class="cell left">
      <img src="../assets/logo-wide-opt.svg" alt="" class="logo">
    </div>
    <div class="cell right"></div>
  </div>
</template>
<script>
export default {
  name: 'SLHeaderBar',
  data() {
    return {
    };
  },
};
</script>
<style scoped>
.header {
  display: grid;
  grid-template-columns: 1fr 1fr;
}

.left {
  justify-self: start;
}

.right {
  justify-self: right;
}

.logo {
  height: 100%;
  max-height: 100%;
}

</style>
