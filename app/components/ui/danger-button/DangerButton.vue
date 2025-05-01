<template>
  <div class="dangerous-button-container">
    <Button
      @click="handleClick"
      :class="{
        'bg-yellow-700': isWaitingForConfirmation,
      }"
      :disabled="isDisabled"
      :variant="$props.variant"
      class="flex items-center justify-center w-auto h-auto"
    >
      <!-- Use slots for custom content -->
      <slot v-if="!isWaitingForConfirmation" name="default-content">
        {{ defaultText }}
      </slot>
      <slot v-else name="confirm-content">
        {{ confirmText }}
      </slot>
      <span v-if="isWaitingForConfirmation && showCountdown" class="countdown">
        ({{ countdownDisplay }})
      </span>
    </Button>
  </div>
</template>

<script>
import { Button } from "../button";
export default {
  props: {
    // Text to show when button is in normal state (used as fallback if no slot is provided)
    defaultText: {
      type: String,
      default: "Delete Item",
    },
    // Text to show when waiting for confirmation (used as fallback if no confirm slot is provided)
    confirmText: {
      type: String,
      default: "Click again to confirm",
    },
    // Time window for confirmation in milliseconds
    confirmationTimeout: {
      type: Number,
      default: 3000, // 3 seconds
    },
    // Whether to show countdown timer
    showCountdown: {
      type: Boolean,
      default: true,
    },
    // Whether the button is disabled
    disabled: {
      type: Boolean,
      default: false,
    },
    variant: {
      type: String,
      default: "primary",
    },
  },
  data() {
    return {
      isWaitingForConfirmation: false,
      timeoutId: null,
      remainingTime: 0,
      countdownInterval: null,
    };
  },
  computed: {
    // This is kept for backward compatibility but no longer used directly in the template
    buttonText() {
      return this.isWaitingForConfirmation
        ? this.confirmText
        : this.defaultText;
    },
    countdownDisplay() {
      return Math.ceil(this.remainingTime / 1000) + "s";
    },
    isDisabled() {
      return this.disabled;
    },
  },
  methods: {
    handleClick() {
      if (!this.isWaitingForConfirmation) {
        // First click - enter confirmation mode
        this.isWaitingForConfirmation = true;
        this.remainingTime = this.confirmationTimeout;

        // Start the timeout to revert to normal state
        this.timeoutId = setTimeout(() => {
          this.resetState();
        }, this.confirmationTimeout);

        // Start countdown timer if enabled
        if (this.showCountdown) {
          this.startCountdown();
        }

        // Emit first-click event
        this.$emit("first-click");
      } else {
        // Second click - confirmation
        this.resetState();

        // Emit confirmed event
        this.$emit("confirmed");
      }
    },
    startCountdown() {
      // Update countdown every 100ms
      this.countdownInterval = setInterval(() => {
        this.remainingTime -= 100;
        if (this.remainingTime <= 0) {
          clearInterval(this.countdownInterval);
        }
      }, 100);
    },
    resetState() {
      this.isWaitingForConfirmation = false;

      if (this.timeoutId) {
        clearTimeout(this.timeoutId);
        this.timeoutId = null;
      }

      if (this.countdownInterval) {
        clearInterval(this.countdownInterval);
        this.countdownInterval = null;
      }
    },
  },
  beforeUnmount() {
    // Clean up any timers when component is destroyed
    this.resetState();
  },
};
</script>

<style>
.confirm-mode {
  @apply bg-yellow-400;
}
</style>
