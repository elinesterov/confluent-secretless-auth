package examples;

public class App {
    public static void main(final String[] args) throws Exception {
        if (args.length < 3) {
          System.out.println("Please provide command line arguments: role configPath topic");
          System.exit(1);
        }

        final String role = args[0];
        final String configPath = args[1];
        final String topic = args[2];

        if ("consumer".equals(role)) {
          // Run the Consumer
          System.out.println("Running consumer");
          ConsumerExample.runConsumer(configPath, topic);
        } else if ("producer".equals(role)) {
          // Run the Producer
          System.out.println("Running producer");
          ProducerExample.runProducer(configPath, topic);
        } else {
          System.out.println("Unknown role: " + role);
          System.exit(1);
        }
    }
}