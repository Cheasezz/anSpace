<script setup lang="ts">
import styles from './styles.module.css'
import ParticlesBg from '../../particlesBg/particlesBg.vue'
</script>

<template>
  <main :class="[styles.main]">
    <picture :class="styles.picture">
      <source
        media="(min-width: 1366px)"
        srcset="../../assets/png/spaceBg3840.png"
      />
      <source
        media="(min-width: 320px)"
        srcset="../../assets/png/spaceBg1920.png"
      />
      <img
        src="../../assets/png/spaceBg3840.png"
        alt="SpaceBg"
      />
    </picture>
    <ParticlesBg />

    <slot></slot>
  </main>
</template>
