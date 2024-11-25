module register_reset #(
    parameter N = 3
) (
    input  logic [N-1:0] data_in,
    input  logic         write, reset,
    output logic [N-1:0] data_out,
);

    logic [N-1:0] data;

    // always_ff @(posedge write, posedge reset) begin
    always_latch begin
        if (reset)
            data <= 0;
        else if (write)
            data <= data_in;
    end

    always_comb begin
        data_out = data;
    end

endmodule