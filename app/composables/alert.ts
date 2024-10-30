import { ref, watch } from "vue";

export const useAlert = () => {
  const pendingAlerts = ref<Alert[]>([]);
  const activeAlert = ref<Alert>();

  type Alert = {
    title: string;
    message: string;
    duration?: number;
    variant?: "destructive" | "default";
  };

  watch(activeAlert, (newAlert) => {
    if (newAlert) {
      const duration = newAlert.duration ?? 5000;
      setTimeout(() => {
        removeAlert();
        setTimeout(() => {
          activateNextAlert();
        }, 300); // This duration is the same as transition on component to make it smooth.
      }, duration);
    }
  });

  function createAlert(alert: Alert) {
    pendingAlerts.value.push(alert);

    if (!activeAlert.value) {
      activateNextAlert();
    }
  }

  function activateNextAlert() {
    if (pendingAlerts.value.length > 0) {
      activeAlert.value = pendingAlerts.value.shift();
    }
  }

  function removeAlert() {
    activeAlert.value = undefined;
  }

  return { createAlert, activeAlert };
};
