<script setup lang="ts">
const pb = usePocketbase();

type BookFacts = {
  rentals: number;
  book_instances: number;
  books: number;
};

type Stats = {
  books: BookFacts;
};

const stats = ref<Stats>({
  books: {
    book_instances: 0,
    books: 0,
    rentals: 0,
  },
});

async function getHomepageStats(): Promise<void> {
  try {
    const data: BookFacts = await pb.send("/homepage-stats", {});
    stats.value.books = { ...data };
  } catch (err) {
    console.log(err);
  }
}

onMounted(async () => {
  await getHomepageStats();
});
</script>
<template>
  <div class="h-screen w-screen flex flex-col items-start justify-start p-5">
    <div class="flex items-center justify-between w-full">
      <img src="/images/logos/massolit-logo.png" class="w-24" alt="" />
    </div>
    <div class="w-full flex flex-col gap-3 mt-20">
      <h1>Services</h1>
      <div class="flex w-full flex-col gap-2">
        <Card class="p-4 flex flex-col gap-4">
          <div class="flex gap-4 items-center">
            <ServiceBookTrackerIcon />
            <h2>Book Tracker</h2>
          </div>

          <div class="flex w-full justify-center gap-4 text-[gray]">
            <div class="flex flex-col gap items-center justify-center">
              <span
                v-if="stats.books.books !== 0"
                class="text-3xl font-bold text-foreground"
                >{{ stats.books.books }}</span
              >
              Books
            </div>

            <div class="flex flex-col gap items-center justify-center">
              <span
                v-if="stats.books.book_instances !== 0"
                class="text-3xl font-bold text-foreground"
                >{{ stats.books.book_instances }}</span
              >
              Copies
            </div>

            <div class="flex flex-col gap items-center justify-center">
              <span
                v-if="stats.books.rentals !== 0"
                class="text-3xl font-bold text-foreground"
                >{{ stats.books.rentals }}</span
              >
              Rented
            </div>
          </div>

          <div>
            <NuxtLink to="/app">
              <Button class="w-full">See all books</Button>
            </NuxtLink>
          </div>
        </Card>

        <Card class="flex flex-col w-full p-4 gap-4 items-center">
          <div class="flex gap-4 w-full items-center">
            <ServiceDetentionsIcon />
            <h2>Detention Tracker</h2>
          </div>

          <div class="w-full">
            <NuxtLink to="/behavior">
              <Button class="w-full">See all detentions</Button>
            </NuxtLink>
          </div>
        </Card>
      </div>
    </div>

    <div></div>
  </div>
</template>
