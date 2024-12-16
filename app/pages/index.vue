<script setup lang="ts">
const pb = usePocketbase();

type BookInfo = {
  rentals: number;
  book_instances: number;
  books: number;
};

type Detentions = {
  pending: number;
};

type Stats = {
  library: BookInfo;
  detentions: Detentions;
};

const stats = ref<Stats>({
  library: {
    book_instances: 0,
    books: 0,
    rentals: 0,
  },
  detentions: {
    pending: 0,
  },
});

async function getHomepageStats(): Promise<void> {
  try {
    const data: Stats = await pb.send("/homepage-stats", {});
    stats.value = { ...data };
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
    <div class="flex items-center w-full justify-between">
      <NuxtLink to="/">
        <img
          src="../public/images/logos/massolit-logo.png"
          class="w-24"
          alt=""
        />
      </NuxtLink>
      <p class="text-sm">{{ pb.authStore.model?.email }}</p>
    </div>

    <div class="w-full lg:w-auto flex flex-col gap-3 mt-20">
      <!-- <h1>Services</h1> -->
      <div class="flex w-full flex-col lg:flex-row gap-2">
        <Card class="p-4 flex flex-col gap-4 lg:gap-8">
          <div class="flex gap-4 items-center">
            <ServiceBookTrackerIcon />
            <h2>Book Tracker</h2>
          </div>

          <div class="flex w-full justify-center gap-4 text-[gray]">
            <div class="flex flex-col gap items-center justify-center">
              <span
                v-if="stats.library.books !== 0"
                class="text-3xl font-bold text-foreground"
                >{{ stats.library.books }}</span
              >
              Books
            </div>

            <div class="flex flex-col gap items-center justify-center">
              <span
                v-if="stats.library.book_instances !== 0"
                class="text-3xl font-bold text-foreground"
                >{{ stats.library.book_instances }}</span
              >
              Copies
            </div>

            <div class="flex flex-col gap items-center justify-center">
              <span
                v-if="stats.library.rentals !== 0"
                class="text-3xl font-bold text-foreground"
                >{{ stats.library.rentals }}</span
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

        <Card
          class="flex flex-col w-full p-4 gap-4 justify-between items-center"
        >
          <div class="flex gap-4 w-full items-center">
            <ServiceDetentionsIcon />
            <h2>Detention Tracker</h2>
          </div>
          <div class="flex w-full justify-center gap-4 text-[gray]">
            <div class="flex flex-col gap items-center justify-center">
              <span
                v-if="stats.detentions.pending !== 0"
                class="text-3xl font-bold text-foreground"
                >{{
                  stats.detentions.pending ? stats.detentions.pending : 0
                }}</span
              >
              {{ stats.detentions.pending === 1 ? "Detention" : "Dententions" }}
            </div>
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
