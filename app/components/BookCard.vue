<template>
  <Card class="p-8 flex flex-col items-start gap-2">
    <h2 class="text-lg font-bold">
      {{ props.book?.title }}
    </h2>

    <div class="flex text-xs justify-center gap-2 items-start">
      <p class="bg-accent rounded px-2 py-0.5">
        {{ props.book?.isbn }}
      </p>
      <p class="bg-primary rounded px-2 py-0.5">
        {{ props.book.id }}
      </p>
    </div>
    <img
      class="w-1/2 place-self-center mt-4"
      :src="props.book.cover_url"
      :alt="'book-cover-' + props.book.title"
    />
  </Card>
</template>

<script setup lang="ts">
interface Props {
  book: {
    id: string;
    isbn: string;
    title: string;
    cover_url: string;
  };
}

const props = defineProps<Props>();
</script>
