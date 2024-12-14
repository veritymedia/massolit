<template>
  <Icon :name="props.name" />
</template>

<script setup lang="ts">
interface Props {
  name: string;
}

const props = withDefaults(defineProps<Props>(), { name: "" });
</script>
